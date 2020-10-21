import { IProduct } from 'shopify-api-node';
import { shopify } from './shopify';

interface Parameters {
  limit: number;
  page_info: string;
}

interface Pagination {
  nextPageParameters: Parameters;
}

export async function listProducts(handler: (products: IProduct[]) => void): Promise<void> {
  let params: Parameters | undefined;
  do {
    const products = await shopify.product.list(params) as IProduct[] & Pagination;
    handler(products);
    params = products.nextPageParameters;
  } while (params !== undefined);
}
