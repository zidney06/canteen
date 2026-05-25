import z from "zod";

export const errorsSchema = z
  .object({
    password: z.array(z.string()),
    selectedKasir: z.array(z.string()),
  })
  .partial();

export type ErrorsType = z.infer<typeof errorsSchema>;

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
