<!--
Di halaman ini, client akan melihat clientkey dan client secret asli mereka.
saat mereka keluar dari halaman ini atau setelah jeda beberapa menit atau detik, halaman akan direfresh
dan client harus segera meyimpan clientkey dan client secret mereka
-->
<script lang="ts">
  import { fetchClientCredetials } from "$lib/repo/repository.js";
  import {
    Card,
    InputGroup,
    Icon,
    CardHeader,
    CardTitle,
  } from "@sveltestrap/sveltestrap";
  import { onMount } from "svelte";

  // req ke BE berdasarkan token url
  // jika OK maka tampilkan client key dan secret dan bisa di copy.
  // jika tidak maka tampilkan pesan "Token tidak valid atau sudah kadaluarsa."

  let clientKey: string = $state("");
  let clientSecret: string = $state("");
  let refreshInterval: number = $state(5);
  let isResOk: boolean = $state(false);
  let isCKCopied = $state(false);
  let isCSCopied = $state(false);

  let { params } = $props();

  $effect(() => {
    // ambil data
    // karena functop $effect tidak boeh berupa promise jadi pakai IIFE async.
    (async () => {
      const res = await fetchClientCredetials(params.token);

      if (!res.isSuccess) {
        isResOk = false;
        return;
      }

      clientKey = res.clientKey;
      clientSecret = res.clientSecret;
      isResOk = true;
    })();

    const intervalId = setInterval(
      async () => {
        const res = await fetchClientCredetials(params.token);

        if (!res.isSuccess) {
          isResOk = false;
          return;
        }

        clientKey = res.clientKey;
        clientSecret = res.clientSecret;
        isResOk = true;
      },
      refreshInterval * 60 * 1000,
    );

    // Fungsi pembersih (cleanup) Svelte 5 tetap aman dan tidak terganggu Promise
    return () => {
      clearInterval(intervalId);
    };
  });

  const handleCopyClientKey = () => {
    navigator.clipboard.writeText(clientKey).then(() => {
      isCKCopied = true;

      setTimeout(() => {
        isCKCopied = false;
      }, 10000);
    });
  };

  const handleCopyClientSecret = () => {
    navigator.clipboard.writeText(clientSecret).then(() => {
      isCSCopied = true;

      setTimeout(() => {
        isCSCopied = false;
      }, 10000);
    });
  };
</script>

<main
  class="d-flex justify-content-center align-items-center"
  style="height: 100vh;"
>
  {#if isResOk}
    <Card
      style="width: 600px;"
      class="p-2 border border-3 border-secondary"
      outline={false}
    >
      <h5 class="text-center">Segera salin dalam {refreshInterval} menit</h5>
      <InputGroup>
        <p class="m-0 align-middle">Client Key: {clientKey}</p>
        <button onclick={handleCopyClientKey} class="btn">
          {#if isCKCopied}
            <Icon name="check-lg" />
          {:else}
            <Icon name="copy" />
          {/if}</button
        >
      </InputGroup>
      <InputGroup>
        <p class="m-0 align-middle">Client Secret: {clientSecret}</p>
        <button onclick={handleCopyClientSecret} class="btn"
          >{#if isCSCopied}
            <Icon name="check-lg" />
          {:else}
            <Icon name="copy" />
          {/if}</button
        >
      </InputGroup>
    </Card>
  {:else}
    <h4>Maaf, token tidak valid atau sudah kadaluarsa.</h4>
  {/if}
</main>

<!--
- pikirkan desain
- mulai buat UIo
- buat logika
-->

<style></style>
