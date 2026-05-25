<script lang="ts">
  import { dummyResponse } from "$lib";
  import SectionLeft from "$lib/components/dashboard/SectionLeft.svelte";
  import SectionRight from "$lib/components/dashboard/SectionRight.svelte";
  import { rawProductsSchema, type Product } from "$lib/types/types";
  import { Container } from "@sveltestrap/sveltestrap";
  import { onMount } from "svelte";

  let dataProduct: Product[] = $state([]);
  let selectedProducts: Product[] = $state([]);
  let totalAmount: number = $state(0);

  onMount(() => {
    // ambil dataProduct daftar barang yang dijual dari backend. jika terdapat error 401 lempar ke login
    const result = rawProductsSchema.safeParse(dummyResponse);

    if (!result.success) {
      console.error(result.error);
      dataProduct = [];
    } else {
      // tambahkan proeprty quantity ke dataProduct agar lebih mudah untuk mengelola UI
      dataProduct = result.data.map((item) => {
        return {
          ...item,
          quantity: 0,
        };
      });
    }
  });

  const addSelectedProducts = (id: string, name: string, price: number) => {
    totalAmount += price;

    const existingProduct = selectedProducts.find((item) => item.id === id);

    if (!existingProduct) {
      selectedProducts.push({ id, name, price, quantity: 1 });
    } else {
      existingProduct.quantity += 1;
    }

    // ubah dataProduct juga agau UI ikut berubah
    dataProduct = dataProduct.map((item) => {
      if (item.id === id) {
        return {
          ...item,
          quantity: item.quantity + 1,
        };
      }
      return item;
    });
  };

  const deleteSelectedProduct = (id: string, quantity: number) => {
    if (quantity <= 0) {
      return;
    }

    totalAmount -= selectedProducts.find((item) => item.id === id)?.price || 0;

    dataProduct = dataProduct.map((product) => {
      if (product.id === id) {
        return { ...product, quantity: product.quantity - 1 };
      }
      return product;
    });

    selectedProducts = selectedProducts.flatMap((item) => {
      if (item.id === id) {
        const newQuantity = item.quantity - 1;

        if (newQuantity > 0) {
          return [{ ...item, quantity: newQuantity }];
        } else {
          return [];
        }
      }

      return [item];
    });
  };

  $inspect(dataProduct);
  $inspect(selectedProducts);
</script>

<Container class="border border-3 rounded" fluid>
  <h3>Dashboard</h3>
  <div class="row">
    <SectionLeft
      {dataProduct}
      {addSelectedProducts}
      {selectedProducts}
      {totalAmount}
      {deleteSelectedProduct}
    />
    <SectionRight {selectedProducts} {totalAmount} />
  </div>
</Container>
