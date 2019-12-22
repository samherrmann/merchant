# shopctl
CLI to manage Shopify stores

## Environment Variables
```
SHOPIFY_SHOP_NAME=shop-name
SHOPIFY_API_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
SHOPIFY_PASSWORD=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
```

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