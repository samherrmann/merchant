export interface ProductConfig {
  /**
   * Type of products.
   */
  type: string,
  /**
   * Path to the product data file, relative to this file.
   */
  dataPath: string;
  /**
   * Product specifications.
   */
  specifications: {
    /**
     * The name of the key/column in the product data file.
     */
    key: string;
    /**
     * The user-visible label for the specification.
     */
    label: string;
    /**
     * The unit of measurement.
     */
    units?: string;
  }[],
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
    key: string;
    valueIndices: number[];
    filenamePattern: string;
    path: string;
  }
}