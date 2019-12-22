import { ProductVariantWeightUnit } from 'shopify-api-node';

export type Product = {
  shopifyId: string;
  vendor: string;
  weight: string;
  weightUnit?: ProductVariantWeightUnit;
} & {
  [specification: string]: string;
}