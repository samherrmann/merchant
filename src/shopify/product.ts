import { IProductVariant, ICreateObjectMetafield } from 'shopify-api-node';
import { DeepPartial } from 'ts-essentials';

export interface NewProduct {
  title: string;
  product_type: string;
  vendor: string;
  tags?: string;
  body_html: string;
  options: { name: string }[];
  variants: NewProductVariant[];
  images: {
    attachment: string;
    filename: string;
  }[];
}

export type NewProductVariant = DeepPartial<IProductVariant> & {
  metafields: ICreateObjectMetafield[];
}