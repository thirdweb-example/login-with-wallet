## Login with Wallet Next Flask

Here we implement the login with wallet flow using the [Next.js](https://nextjs.org/) and [Flask](https://flask.palletsprojects.org/) frameworks.
## Setup

To run the example, first clone this repository, and then `cd` into the `login-with-wallet-next-fiber` directory.

```bash
cd login-with-wallet-next-flask
```

Then run one of the following commands to install the JS dependencies into the `web` folder:

```bash
cd web && npm install && cd ..
# or
cd web && yarn install && cd ..
```

We also need to install the Python dependencies for our backend, which we can do with the following command:

```bash
cd server && pip3 install -r requirements.txt && cd ..
```

Next, you need to create a `.env` file in the `server` folder and add the `ADMIN_PRIVATE_KEY` variable to it (similar to how it is in the `/server/.env.example` file) with the private key of the wallet you want to use as the admin wallet to generate and verify payloads. Your file should use something like the following:

```/server/.env
ADMIN_PRIVATE_KEY=...
```

Now, we need to run the frontend and the backend separately.

We can startup our frontend by running the following command in the root of the `login-with-wallet-next-flask` directory:

```bash
make web
```

Finally, we can startup the backend in a separate terminal window by running the following command:

```bash
make server
```

Now, the demo should be ready to use - and we can navigate to [`http://localhost:3000`](http://localhost:3000) to try it out.