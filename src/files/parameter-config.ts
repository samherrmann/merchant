export interface ParameterConfig {
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
}