import { ProductVariantWeightUnit } from 'shopify-api-node';

export type Product = {
  shopify_id: string;
  vendor: string;
  weight: string;
  weight_unit?: ProductVariantWeightUnit;
} & {
  [specification: string]: string;
}