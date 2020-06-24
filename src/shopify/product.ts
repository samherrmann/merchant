import { IProductVariant, ICreateObjectMetafield } from 'shopify-api-node';

type RecursivePartial<T> = {
  [P in keyof T]?: RecursivePartial<T[P]>;
};

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

export type NewProductVariant = RecursivePartial<IProductVariant> & {
  metafields: ICreateObjectMetafield[];
}