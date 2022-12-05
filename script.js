import http from 'k6/http';
import { check, sleep } from 'k6';

// load domain from environment variable
const domain = __ENV.DOMAIN;
const mx = __ENV.MX;


export const options = {
    stages: [
        { duration: '30s', target: 1000 },
        { duration: '1m', target: 100 },
        { duration: '20s', target: 50 },
    ],
};

export default function () {
    const res = http.get('http://app:8080/.well-known/mta-sts.txt');

    check(res, { 'status was 200': (r) => r.status == 200 });
    // check that the response body contains the text with valid MTA-STS policy version
    check(res, { 'response body contains valid MTA-STS policy version': (r) => r.body.includes('version: STSv1') });
    // check that the response body contains the text with valid MTA-STS policy mode
    check(res, { 'response body contains valid MTA-STS policy mode': (r) => r.body.includes('mode: testing') });
    // check that the response body contains the text with valid MTA-STS policy max_age
    check(res, { 'response body contains valid MTA-STS policy max_age': (r) => r.body.includes('max_age: 86400') });
    // check that the response body contains the text with valid MTA-STS policy mx
    check(res, {
        'response body contains valid MTA-STS policy mx': (r) => r.body.includes
            (`mx: ${mx}`)
    });
    sleep(1);
}
