import { ProductConfig } from '../files/product-config';
import { Product } from '../files/product';
import { IProduct, ICreateObjectMetafield, ProductVariantWeightUnit } from 'shopify-api-node';
import { maskString } from '../utils/mask-string';
import mustache from 'mustache';
import { shopify } from './shopify';
import { readImages } from '../files/read-image';
import { ParameterConfig } from '../files/parameter-config';
import util from 'util';
import { NewProduct, NewProductVariant } from './product';

// Override Mustache's escape function to not HTML-escape variables.
mustache.escape = (text: string): string => text;

/**
 * Returns specifications as metafields.
 */
function createMetafields(p: Product, c: ProductConfig): ICreateObjectMetafield[] {
  return c.specifications.reduce<ICreateObjectMetafield[]>((prev, curr) => {
    let label = curr.label;
    if (curr.units) {
      label = `${label} [${curr.units}]`;
    }
    prev.push({
      key: label,
      value: p[curr.key],
      /* eslint-disable-next-line @typescript-eslint/camelcase */
      value_type: 'string',
      namespace: 'specifications'
    });
    return prev;
  }, []);
}

/**
 * Returns the title for the provided product.
 */
function createTitle(product: Product, config: ProductConfig): string {
  return mustache.render(config.title, product);
}

/**
 * Returns the specifications of the provided product in an HTML table.
 */
function createSpecificationTable(product: Product, c: ProductConfig): string {
  const rows = c.specifications.reduce<string[]>((prev, curr) => {
    let label = curr.label;
    if (curr.units) {
      label = `${label} [${curr.units}]`
    }
    prev.push(`
    <tr>
      <th>${label}</th>
      <td>${product[curr.key]}</td>
    </tr>
    `);
    return prev;
  }, []);

  return `<table>${rows.join('')}</table>`
}

/**
 * Returns an HTML specification table for each provided variant.
 * If more than one variant is provided, then each table is wrapped
 * in a `<template>` element with an element ID of the form
 * `specification-option1-option2-option3-template`.
 */
function createSpecificationTables(variants: Product[], config: ProductConfig): string[] {
  return variants.map(v => {
    const table = createSpecificationTable(v, config);
    const options = [];
    if (config.option1) {
      options.push(v[config.option1.key]);
    }
    if (config.option2) {
      options.push(v[config.option2.key]);
    }
    if (config.option3) {
      options.push(v[config.option3.key]);
    }
    return `<template id="specification-${options.join('-')}-template">${table}</template>`
  });
}

/**
 * Returns the image filenames based on the configuration.
 */
function imageFilenames(product: Product, config: ProductConfig): string[] {
  const imageConfig = config.image;
  if (imageConfig && product[imageConfig.key]) {  
    const filenames = product[imageConfig.key].split('|').map(fn => fn.trim());
    return filenames.map(filename => {
      // If the image configuration also defines a `charIndices` and a
      // `filenamePattern` property, then mask the filename value using those
      // properties.
      if (imageConfig.charIndices && imageConfig.filenamePattern) {
        filename = maskString(
          filename,
          imageConfig.charIndices,
          imageConfig.filenamePattern
        );
      }
      let path = `${imageConfig.dir}/${filename}`;
      if (!path.endsWith('.jpg')) {
        path += '.jpg';
      }
      return path;
    });
  }
  return [];
}

/**
 * Returns the custom product property names.
 */
function createOptions(config: ProductConfig): { name: string }[] {
  return [config.option1, config.option2, config.option3]
    .filter((opt): opt is ParameterConfig => !!opt)
    .map(opt => {
      let label = opt.label;
      if (opt.units) { label = `${label} [${opt.units}]`; }
      return { name: label };
    });
}

/**
 * Returns an array of any duplicate variants that may exist within the product.
 * Returns an empty array if no duplicates exist.
 */
function duplicateVariants(product: NewProduct): [NewProductVariant, NewProductVariant][] {
  const variants = product.variants || [];
  const options = variants.map(v => (v?.option1 || '') + (v?.option2 || '') + (v?.option3 || ''));
  return options.reduce<[NewProductVariant, NewProductVariant][]>((acc, curr, i) => {
    const lastIndex = options.lastIndexOf(curr);
    if (i !== lastIndex) {
      acc.push([variants[i], variants[lastIndex]]);
    }
    return acc;
  }, []);
}

/**
 * Create a product in the Shopify store.
 */
export async function createProduct(variants: Product[], config: ProductConfig): Promise<IProduct> {
  const defaultProduct = variants[0];
  const title = createTitle(defaultProduct, config);
  const images = await readImages(imageFilenames(defaultProduct, config));

  /* eslint-disable @typescript-eslint/camelcase */
  const product: NewProduct = {
    title: title,
    vendor: defaultProduct.vendor,
    product_type: config.type,
    body_html: createSpecificationTables(variants, config).join(''),
    options: createOptions(config),
    variants: variants.map(v => {
      const variant: NewProductVariant = {
        sku: v[config.skuKey || 'sku'],
        barcode: v[config.barcodeKey || 'barcode'],
        inventory_management: 'shopify',
        weight: parseFloat(v[config.weightKey || 'weight']),
        weight_unit: v[config.weightUnitKey || 'weight_unit'] as ProductVariantWeightUnit || 'kg',
        metafields: createMetafields(v,config),
        price: v[config.priceKey || 'price']
      };
      if (config.option1) {
        variant.option1 = v[config.option1.key]
      }
      if (config.option2) {
        variant.option2 = v[config.option2.key]
      }
      if (config.option3) {
        variant.option3 = v[config.option3.key]
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

  const dups = duplicateVariants(product);
  if (dups.length > 0) {
    throw util.inspect({ duplicates: dups }, false, null, true);
  }

  return shopify.product.create(product).catch(err => {
    throw util.inspect({ err, product }, false, null, true);
  })
}
