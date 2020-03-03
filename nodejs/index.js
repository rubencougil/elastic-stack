// Add this to the VERY top of the first file loaded in your app
const fetch = require('node-fetch');
const fibonacci = require ('fibonacci');
const apm = require('elastic-apm-node').start({
  // Override service name from package.json
  // Allowed characters: a-z, A-Z, 0-9, -, _, and space
  serviceName: 'nodejs',

  // Use if APM Server requires a token
  secretToken: '6hndgxqdkq6dcd9rfl5ds54g',

  // Set custom APM Server URL (default: http://localhost:8200)
  serverUrl: 'http://localhost:8200',
});

const Tracer = require('elastic-apm-node-opentracing');

// // Pass the Elastic APM agent as an argument to the OpenTracing tracer
const tracer = new Tracer(apm);

const app = require('express')();

app.get('/', function (req, res) {
    const span = tracer.startSpan('my-first-span');
    // ... do some work ...
    fetch('http://localhost:8080/')
        .then(res => res.text())
        .then(body => console.log(body));
    fetch('http://localhost:8080/create', { method: 'POST', body: {} })
        .then(res => res.text())
        .then(body => console.log(body));
    span.finish()
  res.send('Hello World!')
});

app.get('/expensive_calc', function (req, res) {
    const span = tracer.startSpan('expensive calc');
    const bigNumber = fibonacci.iterate (10000);
    res.send(bigNumber);
    span.finish()
})

app.listen(3200)