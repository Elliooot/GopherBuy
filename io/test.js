import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();

// Load proto in init context
client.load(['/src/api'], 'order.proto');

export const options = {
    stages: [
        { duration: '30s', target: 100 }, // Concurrent number increaces to 100 within 30 secs
        { duration: '1m', target: 500 }, // keep 500 concurrent number last for 1 minute
        { duration: '30s', target: 0 }, // cooldown to 0 connection
    ],
    thresholds: {
        'grpc_req_duration': ['p(95)<500', 'p(99)<1000'], // 95% < 500ms, 99% < 1s
        'checks': ['rate>0.9'], // 90% of checks should pass
    },
};

export default function () {
    if (!client.isConnected) {
        try {
            client.connect('app:50051', {
                plaintext: true,
                timeout: '10s',
            });
            console.log('Connected to gRPC server');
        } catch (error) {
            console.error('Failed to connect:', error);
            return;
        }
    }

    const data = {
        user_id: Math.floor(Math.random() * 1000000),
        product_id: 101,
        promo_id: 2,
        quantity: 1,
    };

    try {
        // Send RPC request
        const response = client.invoke('api.OrderService/CreateFlashOrder', data);

        check(response, {
            'status is OK': (r) => r && r.status === grpc.StatusOK,
            'response exists': (r) => r && r.message,
            'resp status is 202': (r) => r && r.message.status === 202, // response code when in queue
            'has order_sn': (r) => r && r.message && r.message.orderSn,
        });

        if (response.status !== grpc.StatusOK) {
            console.error('Request failed:', JSON.stringify(response));
        }

    } catch (error) {
        console.error('Request error:', error);
    }
    
    sleep(0.1); // The interval for user requests is 0.1 seconds.
}

export function teardown() {
    if (client.isConnected) {
        client.close();
        console.log('Disconnected from gRPC server');
    }
}