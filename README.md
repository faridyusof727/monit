### Uptime Agent

To run locally

1. Install Go. 
2. Clone this repo.
3. Run `go mod tidy` and `go mod vendor`.
4. Create your `.env` file. Refer to `.env.example`.
5. Add `firebase-admin-sdk.json` in your root folder. This is your Firebase Service Account key. Refer [here](https://firebase.google.com/docs/admin/setup#initialize-sdk).
6. Create a counter file in `store/counter/counter.json` like below:

```json
{
  "Current": 1, 
  "Max": 2000 // this is unused stuff. still not implemented.
}
```

7. Give it a go by `go run .` or if you want to build `go build .`.


