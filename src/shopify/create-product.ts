import { ProductConfig } from '../files/product-config';
import { Product } from '../files/product';
import Shopify, { IProductImage, IProduct } from 'shopify-api-node';
import { maskString } from '../utils/mask-string';
import { existsSync, readFileSync } from 'fs-extra';
import mustache from 'mustache';
import { shopify } from './shopify';

// Override Mustache's escape function to not HTML-escape variables.
mustache.escape = text => text;

/**
 * Create a product in the Shopify store.
 */
export function createProduct(p: Product, c: ProductConfig): Promise<IProduct> {
  const title = createTitle(p, c);
  const table = createSpecificationsTable(p, c);
  const image = readImage(p, c);

  const product: RecursivePartial<Shopify.IProduct> = {
    title: title,
    body_html: table,
    vendor: p.vendor,
    product_type: c.type,
    tags: [p.vendor, p.name].join(', '),
    variants: [{
      inventory_management: 'shopify',
      weight: parseFloat(p.weight),
      weight_unit: p.weightUnit || 'kg'
    }],
    // IProductImage interface does not currently define the attachment
    // property as defined by the Shopify documentation:
    // https://help.shopify.com/en/api/reference/products/product#create-2019-10
    images: [
      ...image ? [{ attachment: image }] : []
    ] as Partial<IProductImage & { attachment: string }>[]
  };
  return shopify.product.create(product);
}

/**
 * Returns the specifications of the provided product in an HTML table.
 */
function createSpecificationsTable(p: Product, c: ProductConfig): string {
  const rows = c.specifications.reduce<string[]>((prev, curr) => {
    let label = curr.label;
    if (curr.units) {
      label = `${label} [${curr.units}]`
    }
    prev.push(`
    <tr>
      <th>${label}</th>
      <td>${p[curr.key]}</td>
    </tr>
    `);
    return prev;
  }, []);

  return `<table>${rows.join('')}</table>`
};

/**
 * Returns the title for the provided product.
 */
function createTitle(p: Product, c: ProductConfig): string {
  return mustache.render(c.title, p);
}

function readImage(p: Product, c: ProductConfig): string | undefined {
  const image = c.image;
  if (!(image && p[image.key])) { return; }

  const filename = maskString(
    p[image.key],
    image.valueIndices,
    image.filenamePattern
  );

  const filePath = `${image.path}/${filename}.jpg`;
  if (!existsSync(filePath)) { return; }

  return readFileSync(filePath).toString('base64');
}

type RecursivePartial<T> = {
  [P in keyof T]?: RecursivePartial<T[P]>;
};