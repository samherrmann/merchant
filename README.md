# shopctl
CLI to manage Shopify stores.

## Installation
```
npm install -g shopctl
```

## Getting Started

1. Set the required environment variables:
    ```
    SHOPIFY_SHOP_NAME=my-shop-name
    SHOPIFY_API_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
    SHOPIFY_PASSWORD=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
    ```
    Environment variables may be defined in a `.env` file in the root directory
    of your project.

1. In the root directory of your project, add a file named `shopctl.json`
   with content of the following form:

    `project-root/shopctl.json`:
    ```json
    {
      "products": [
        "products/shoes.json",
        "products/hats.json",
        "products/shirts.json"
      ]
    }
    ```

1. Add each product type definition file listed in `shopctl.json` in the form
   as follows:

    `project-root/products/shoes.json`:
    ```json
    {
      "type": "Shoe",
      "dataPath": "shoes.csv",
      "specifications": [{
        "key": "size",
        "label": "Size"
      }, {
        "key": "color",
        "label": "Color"
      }],
      "title": "{{name}}, {{size}}, {{color}}"
    }
    ```
    * `dataPath` defines the path of the products data file (see next step)
      relative to this file.
    * `shopctl` uses [Mustache](http://mustache.github.io/) to generate the
      product title. The `title` string defined in the product type definition
      file may therefore be written using [Mustache template
      syntax](http://mustache.github.io/mustache.5.html). `{{ tags }}` in the
      template may reference the column names defined in the product data file.

1. Add the products data file for each product type:

    `project-root/products/shoes.csv`:
    ```csv
    shopifyId,brand,name, size, color
    ,33969,Asics, DynaFlyte 4, 10, Blue
    ,33983,Adidas,Ultaboost 20, 11, Red
    ,40127,New Balance,Fresh Foam Beacon V2, 12, Grey
    ```
    The `shopifyId` column may be left blank and will be populated by `shopctl`
    once the product is added to the store.

The project is now all set up and ready to be used by `shopctl` from the project root directory:

```
shopctl help
shopctl products count
shopctl products push
```

## Development

1. Build and watch:

    ```sh
    npm start
    ```

2. In a separate terminal, run the CLI:
    ```sh
    npm link
    shopctl help
    ```