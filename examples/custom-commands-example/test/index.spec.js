const request = require('supertest');
const expect = require('expect');
const app = require('../index.js').app;

describe("get al users",() =>{
    it('should return an array with one test user', (done)=>{
        request(app)
        .get('/users')
        .expect('Content-Type', /json/)
        .expect(200)
        .expect(response => {
            //console.log(response)
            expect(response.body.length).toBe(1);
            expect(response.body[0].username).toBe('test');
        })
        .end(done);
    })
});

describe("post one user",() =>{
    it('should return the brand new just created user', (done)=>{

        let user = {username: "juanzy", age:23};

        request(app)
        .post('/users')
        .send(user)
        .expect('Content-Type', /json/)
        .expect(200)
        .expect(response => {
            console.log(response.body)
            expect(response.body.id && response.body.username && response.body.age).toBeTruthy();
            expect(response.body.username).toBe(user.username);
            expect(response.body.id).toExist();
        })
        .end(done);
    })
});
