import { ISmartCollection } from 'shopify-api-node';
import { shopify } from './shopify';

export async function listCollections(): Promise<ISmartCollection[]> {
  return shopify.smartCollection.list()
}
