import { ParameterConfig } from './parameter-config';

export interface ProductConfig {
  /**
   * Type of products.
   */
  type: string;
  /**
   * Path to the product data file, relative to the working directory.
   */
  dataPath: string;
  /**
   * Product specifications.
   */
  specifications: ParameterConfig[];
  /**
   * The product title. The title may be constructed from product data using
   * the following form:
   *
   * ```
   * {{ name }}, {{ dataFieldA }}, {{ dataFieldB}}
   * ```
   */
  title: string;
  /**
   * Product image.
   */
  image?: {
    /**
     * Key of the product property from which to construct the filename.
     */
    key: string;
    /**
     * The indices of the characters to extract from the product property.
     * The extracted characters are used to generate the filename.
     */
    charIndices: number[];
    /**
     * Pattern defining the format of the filename. All `#` characters in the
     * pattern are replaced with the extracted characters per `charIndices`, in
     * sequential order. The number of `#` characters defined in the pattern
     * should match the length of the `charIndices` array.
     */
    filenamePattern: string;
    /**
     * Path to the directory containing the image, relative to the working directory.
     */
    dir: string;
  };

  weightKey?: string;

  skuKey: string;

  barcodeKey: string;

  option1?: ParameterConfig;

  option2?: ParameterConfig;

  option3?: ParameterConfig;
}