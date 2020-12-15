# hackathon
Submission to the 2020 Textkernel Hackathon

This is a Proof of Concept implementation of a Routing Slip Pattern for
Kubernetes. The goal is to create shared libraries for the heavy lifting, so
each Component can focus on it's own implementation.

For this the following `golang` modules were created:
* [`gobernate`](https://github.com/SebastiaanPasterkamp/gobernate) - a simple
  HTTP mux Router wrapper providing the bare minimum for a Kubernetes pod.
  Besides graceful shutdown the following routes are available out-of-the-box:
  * `GET /version` to return a JSON structure with the service.
  * `GET /health` to return http.StatusOK to indicate the (web) service is
    running.
  * `GET /readiness` to signal when the service is ready to serve.
  * `GET /metrics` to return Prometheus formatted metric data.   
* [`gony express`](https://github.com/SebastiaanPasterkamp/gonyexpress) - an
  AMQP routing library doing the heavy lifting of handing incoming and outgoing
  messages. All a Component needs to provide is a call-back function to be
  called for every message.

With these 2 packages the following example Components were created:
* `postoffice` - instead of expecting every Producer to know all route
  configurations, one simply sends the message to the `postoffice` first. The
  routing slip matching the name is attached.
* `attach` - attaches a file by name. Demonstrates build-in base64 encoding
  for Documents.
* `checksum` - computes the md5 on the attached file. Demonstrates streaming
  reading/writing into Documents.
* `delay` - a simple sleep component demonstrating the usage of Step arguments.
* `unstable` - a component that randomly fails, to demonstrate the error
  handling capabilities.

And the `producer` sends a number of initial messages to the post office, with
random routes and filenames to be attached.

To launch the setup run:

```bash
docker-compose build
docker-compose up
```

and to send more messages repeat the following:

```bash
docker-compose up producer
```
