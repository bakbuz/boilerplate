import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();
client.load(['../../../api/protos'], 'catalog/v1/brand.proto');

export const options = {
    stages: [
        { duration: '10s', target: 10 },
        { duration: '20s', target: 10 },
        { duration: '5s', target: 0 },
    ],
};

export default () => {
    client.connect('127.0.0.1:50051', {
        plaintext: true,
    });

    const data = {};
    const response = client.invoke('catalog.v1.BrandService/List', data);

    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });

    client.close();
    sleep(1);
};
