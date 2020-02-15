import Shopify from 'shopify-api-node';

const envVars = {
  shopName: 'SHOPIFY_SHOP_NAME',
  apiKey: 'SHOPIFY_API_KEY',
  password: 'SHOPIFY_PASSWORD'
};

const shopName = process.env[envVars.shopName];
if (!shopName) {
  console.error(`${envVars.shopName} is not defined.`);
  process.exit(1);
}

const apiKey = process.env[envVars.apiKey];
if (!apiKey) {
  console.error(`${envVars.apiKey} is not defined.`);
  process.exit(1);
}

const password = process.env[envVars.password];
if (!password) {
  console.error(`${envVars.password} is not defined.`)
  process.exit(1);
}

export const shopify = new Shopify({
  autoLimit: true,
  apiVersion: '2019-07',
  shopName: shopName,
  apiKey: apiKey,
  password: password
});