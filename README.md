# shopctl
CLI to manage Shopify stores

## Installation
```
npm install -g shopctl
```

## Environment Variables
```
SHOPIFY_SHOP_NAME=shop-name
SHOPIFY_API_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
SHOPIFY_PASSWORD=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
```
Environment variables may be defined in a `.env` file. `shopctl` looks for the
`.env` file in the current working directory.

## Development

1. Build and watch:

    ```sh
    npm start
    ```

2. In a separate terminal, run the CLI:
    ```sh
    npm link
    shopctl --help
    ```