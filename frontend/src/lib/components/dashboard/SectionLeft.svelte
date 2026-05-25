<script lang="ts">
  import type { Product } from "$lib/types/types";
  import {
    Button,
    ButtonGroup,
    Card,
    CardBody,
    CardText,
  } from "@sveltestrap/sveltestrap";

  interface Props {
    dataProduct: Product[];
    selectedProducts: Product[];
    totalAmount: number;
    addSelectedProducts: (id: string, name: string, price: number) => void;
    deleteSelectedProduct: (id: string, quantity: number) => void;
  }

  const {
    dataProduct,
    addSelectedProducts,
    totalAmount,
    deleteSelectedProduct,
  }: Props = $props();
</script>

<section class="col border overflow-auto p-2" style="height: 90vh;">
  <div class="border p-2">
    <p>Total Amount: {totalAmount}</p>
  </div>
  <div class="border p-2" style="height: 200px;">
    <h4 class="text-center">Daftar produk</h4>
    {#each dataProduct as item}
      <Card class="my-2">
        <CardBody>
          <CardText>
            Nama barang: {item.name}. Harga: {item.price}
          </CardText>
          <ButtonGroup class="align-items-center">
            <Button
              onclick={() =>
                addSelectedProducts(item.id, item.name, item.price)}>+</Button
            >
            <span class="align-middle px-2">{item.quantity}</span>
            <Button
              onclick={() => deleteSelectedProduct(item.id, item.quantity)}
              >-</Button
            >
          </ButtonGroup>
        </CardBody>
      </Card>
    {/each}
  </div>
</section>
