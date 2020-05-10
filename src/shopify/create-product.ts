import { ProductConfig } from '../files/product-config';
import { Product } from '../files/product';
import Shopify, { IProduct, IProductVariant } from 'shopify-api-node';
import { maskString } from '../utils/mask-string';
import mustache from 'mustache';
import { shopify } from './shopify';
import { readImages } from '../files/read-image';
import { ParameterConfig } from '../files/parameter-config';

// Override Mustache's escape function to not HTML-escape variables.
mustache.escape = (text: string): string => text;

/**
 * Returns the title for the provided product.
 */
function createTitle(p: Product, c: ProductConfig): string {
  return mustache.render(c.title, p);
}

/**
 * Returns the image filenames based on the configuration.
 */
function imageFilenames(p: Product, c: ProductConfig): string[] {
  const imageConfig = c.image;
  if (imageConfig && p[imageConfig.key]) {  
    const filenames = p[imageConfig.key].split('|').map(fn => fn.trim());
    return filenames.map(filename => {
      // If the image configuration also defines a `charIndices` and a
      // `filenamePattern` property, then mask the filename value using those
      // properties.
      if (imageConfig.charIndices && imageConfig.filenamePattern) {
        filename = maskString(
          filename,
          imageConfig.charIndices,
          imageConfig.filenamePattern
        ) + '.jpg';
      }
      return `${imageConfig.dir}/${filename}`;
    });
  }
  return [];
}

/**
 * Returns the custom product property names.
 */
function createOptions(c: ProductConfig): { name: string }[] {
  return [c.option1, c.option2, c.option3]
    .filter((opt): opt is ParameterConfig => !!opt)
    .map(opt => {
      let label = opt.label;
      if (opt.units) { label = `${label} [${opt.units}]`; }
      return { name: label };
    });
}

/**
 * Create a product in the Shopify store.
 */
export async function createProduct(variants: Product[], c: ProductConfig): Promise<IProduct> {
  const defaultProduct = variants[0];
  const title = createTitle(defaultProduct, c);
  const images = await readImages(imageFilenames(defaultProduct, c));

  /* eslint-disable @typescript-eslint/camelcase */
  const product: NewProduct = {
    title: title,
    vendor: defaultProduct.vendor,
    product_type: c.type,
    tags: [defaultProduct.vendor, defaultProduct.name].join(', '),
    options: createOptions(c),
    variants: variants.map(v => {
      const variant: NewProductVariant = {
        inventory_management: 'shopify',
        weight: parseFloat(v[c.weightKey || 'weight']),
        weight_unit: v.weight_unit || 'kg'
      };
      if (c.option1) {
        variant.option1 = v[c.option1.key]
      }
      if (c.option2) {
        variant.option2 = v[c.option2.key]
      }
      if (c.option3) {
        variant.option3 = v[c.option3.key]
      }
      return variant;
    }),
    images: images.map(image => {
      return {
        attachment: image.base64,
        filename: image.filename
      };
    })
  };
  /* eslint-enable @typescript-eslint/camelcase */
  return shopify.product.create(product);
}

type RecursivePartial<T> = {
  [P in keyof T]?: RecursivePartial<T[P]>;
};

interface Image {
  base64: string;
  filename: string;
}

type NewProduct = RecursivePartial<Shopify.IProduct & {
  images: {
    attachment: string;
    filename: string;
  }[];
}>;

type NewProductVariant = RecursivePartial<IProductVariant>