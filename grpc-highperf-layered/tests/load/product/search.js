import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();
client.load(['../../../api/protos'], 'catalog/v1/product.proto');

export const options = {
    stages: [
        { duration: '10s', target: 5 },
        { duration: '20s', target: 5 },
        { duration: '5s', target: 0 },
    ],
};

export default () => {
    client.connect('127.0.0.1:50051', {
        plaintext: true,
    });

    const data = {
        name: 'Product', // Common name to get results
        limit: 10,
        // brand_id: 1, // Optional: filter by brand
        // last_seen_id: "", // Optional: for pagination
    };
    const response = client.invoke('catalog.v1.ProductService/Search', data);

    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
        'has items': (r) => r.message && r.message.items && r.message.items.length >= 0, // Should at least return an empty list, not error
    });

    client.close();
    sleep(1);
};
