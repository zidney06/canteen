import z from "zod";

export type ErrorsType = {
  password?: string[];
  selectedKasir?: string[];
};

export const rawProductSchema = z.object({
  id: z.string(),
  name: z.string(),
  price: z.number().positive().max(100000),
});

export const rawProductsSchema = z.array(rawProductSchema);

export type RawProductSchema = z.infer<typeof rawProductSchema>;

export interface Product extends RawProductSchema {
  quantity: number;
}
