# OpenTelemetry Trace Propagation - Go

This repo provides an extremely quick and simple example of one way you can propagate OpenTelemetry trace context across process boundaries with Go. Although this implementation is in a single process, the same technique can be used any time you serialise data into a slice of bytes to pass between processes (E.g. a queue, event stream, etc.)

## How It Works

The `Carrier` implements the `propagator.TextMapCarrier` interface using a `map[string]string` and the required methods. This allows the `TextMapPropagator` we define at start time to inject the Trace and Span information from the context into the `Carrier` and extract it from the `Carrier` into local context when we send it between "processes."

The example in this code is trivial but the same technique can be used to include Trace information in the body of messages, as part of events emitted to an event stream or you could do something [truly wild and awesome](https://twitter.com/MissAmyTobey/status/1437209479477534721)!

## Example Output

Here's some example output showing how `processA` and `processB` share the same parent trace:

```json
{
	"Name": "processB",
	"SpanContext": {
		"TraceID": "af34a4b214932ba9c6a411a9fa5c31f6",
		"SpanID": "07d34acae71269de",
		"TraceFlags": "01",
		"TraceState": "",
		"Remote": false
	},
	"Parent": {
		"TraceID": "af34a4b214932ba9c6a411a9fa5c31f6",
		"SpanID": "15f6a29cc64e1482",
		"TraceFlags": "01",
		"TraceState": "",
		"Remote": true
	},
	"SpanKind": 1,
	"StartTime": "2021-09-28T21:46:08.123697231+10:00",
	"EndTime": "2021-09-28T21:46:08.123698614+10:00",
	"Attributes": null,
	"Events": [
		{
			"Name": "context extracted from carrier, now we're in the right parent trace",
			"Attributes": null,
			"DroppedAttributeCount": 0,
			"Time": "2021-09-28T21:46:08.123698186+10:00"
		}
	],
	"Links": null,
	"Status": {
		"Code": "Unset",
		"Description": ""
	},
	"DroppedAttributes": 0,
	"DroppedEvents": 0,
	"DroppedLinks": 0,
	"ChildSpanCount": 0,
	"Resource": [
		{
			"Key": "service.name",
			"Value": {
				"Type": "STRING",
				"Value": "unknown_service:main"
			}
		},
		{
			"Key": "telemetry.sdk.language",
			"Value": {
				"Type": "STRING",
				"Value": "go"
			}
		},
		{
			"Key": "telemetry.sdk.name",
			"Value": {
				"Type": "STRING",
				"Value": "opentelemetry"
			}
		},
		{
			"Key": "telemetry.sdk.version",
			"Value": {
				"Type": "STRING",
				"Value": "1.0.0"
			}
		}
	],
	"InstrumentationLibrary": {
		"Name": "processB",
		"Version": "",
		"SchemaURL": ""
	}
}
{
	"Name": "processA",
	"SpanContext": {
		"TraceID": "af34a4b214932ba9c6a411a9fa5c31f6",
		"SpanID": "15f6a29cc64e1482",
		"TraceFlags": "01",
		"TraceState": "",
		"Remote": false
	},
	"Parent": {
		"TraceID": "af34a4b214932ba9c6a411a9fa5c31f6",
		"SpanID": "d89a5c8913713a27",
		"TraceFlags": "01",
		"TraceState": "",
		"Remote": false
	},
	"SpanKind": 1,
	"StartTime": "2021-09-28T21:46:08.123626787+10:00",
	"EndTime": "2021-09-28T21:46:08.12370727+10:00",
	"Attributes": null,
	"Events": [
		{
			"Name": "context injected into carrier: {Fields:map[traceparent:00-af34a4b214932ba9c6a411a9fa5c31f6-15f6a29cc64e1482-01]}",
			"Attributes": null,
			"DroppedAttributeCount": 0,
			"Time": "2021-09-28T21:46:08.123641363+10:00"
		},
		{
			"Name": "sending message to processB",
			"Attributes": null,
			"DroppedAttributeCount": 0,
			"Time": "2021-09-28T21:46:08.12367818+10:00"
		},
		{
			"Name": "message sent to processB",
			"Attributes": null,
			"DroppedAttributeCount": 0,
			"Time": "2021-09-28T21:46:08.12370701+10:00"
		}
	],
	"Links": null,
	"Status": {
		"Code": "Unset",
		"Description": ""
	},
	"DroppedAttributes": 0,
	"DroppedEvents": 0,
	"DroppedLinks": 0,
	"ChildSpanCount": 0,
	"Resource": [
		{
			"Key": "service.name",
			"Value": {
				"Type": "STRING",
				"Value": "unknown_service:main"
			}
		},
		{
			"Key": "telemetry.sdk.language",
			"Value": {
				"Type": "STRING",
				"Value": "go"
			}
		},
		{
			"Key": "telemetry.sdk.name",
			"Value": {
				"Type": "STRING",
				"Value": "opentelemetry"
			}
		},
		{
			"Key": "telemetry.sdk.version",
			"Value": {
				"Type": "STRING",
				"Value": "1.0.0"
			}
		}
	],
	"InstrumentationLibrary": {
		"Name": "processA",
		"Version": "",
		"SchemaURL": ""
	}
}
{
	"Name": "main",
	"SpanContext": {
		"TraceID": "af34a4b214932ba9c6a411a9fa5c31f6",
		"SpanID": "d89a5c8913713a27",
		"TraceFlags": "01",
		"TraceState": "",
		"Remote": false
	},
	"Parent": {
		"TraceID": "00000000000000000000000000000000",
		"SpanID": "0000000000000000",
		"TraceFlags": "00",
		"TraceState": "",
		"Remote": false
	},
	"SpanKind": 1,
	"StartTime": "2021-09-28T21:46:08.123616078+10:00",
	"EndTime": "2021-09-28T21:46:08.123709294+10:00",
	"Attributes": null,
	"Events": null,
	"Links": null,
	"Status": {
		"Code": "Unset",
		"Description": ""
	},
	"DroppedAttributes": 0,
	"DroppedEvents": 0,
	"DroppedLinks": 0,
	"ChildSpanCount": 1,
	"Resource": [
		{
			"Key": "service.name",
			"Value": {
				"Type": "STRING",
				"Value": "unknown_service:main"
			}
		},
		{
			"Key": "telemetry.sdk.language",
			"Value": {
				"Type": "STRING",
				"Value": "go"
			}
		},
		{
			"Key": "telemetry.sdk.name",
			"Value": {
				"Type": "STRING",
				"Value": "opentelemetry"
			}
		},
		{
			"Key": "telemetry.sdk.version",
			"Value": {
				"Type": "STRING",
				"Value": "1.0.0"
			}
		}
	],
	"InstrumentationLibrary": {
		"Name": "main",
		"Version": "",
		"SchemaURL": ""
	}
}
```