import http from 'k6/http';
import { check } from 'k6';

export const options = {
    stages: [
        { target: 50, duration: '30s' },
        { target: 50, duration: '1m' },
        { target: 0, duration: '30s' },
    ],
    thresholds: {
        http_req_failed: ['rate<0.01'],
        http_req_duration: ['p(95)<300'],
    },
};

export default function () {
    const res = http.get(`http://localhost:1919/api/todo/lists`);

    check(res, {
        'status is 200': (r) => r.status === 200,
        'response body': (r) => {
            const result = JSON.parse(r.body);
            return result.count && result.count > 0;
        },
    });
}
