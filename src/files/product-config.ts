export interface ProductConfig {
  type: string,
  dataPath: string;
  specifications: {
    key: string;
    label: string;
    units?: string;
  }[],
  title: string;
  image?: {
    key: string;
    valueIndices: number[];
    filenamePattern: string;
    path: string;
  }
}