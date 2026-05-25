<script lang="ts">
  import type { Product } from "$lib/types/types";
  import { Button } from "@sveltestrap/sveltestrap";
  import { Html5Qrcode } from "html5-qrcode";
  import { type QrcodeResult } from "html5-qrcode/esm/core";
  import { onMount } from "svelte";

  interface Props {
    selectedProducts: Product[];
    totalAmount: number;
  }

  interface DataBuyer {
    name: string;
    nis: string;
  }

  const { selectedProducts, totalAmount } = $props();

  let cameraId: string = "";
  let html5QrCode: any;
  let isDataReceived: boolean = $state(false);
  let dataBuyer: DataBuyer = $state({
    name: "",
    nis: "",
  });

  onMount(() => {
    html5QrCode = new Html5Qrcode("reader");
    Html5Qrcode.getCameras()
      .then((devices) => {
        if (devices && devices.length) {
          cameraId = devices[0].id;
          console.log("Successfully to get camera");
        }
      })
      .catch((err) => {
        console.log(err);
      });
  });

  const getBuyerData = async (qrDecoded: string) => {
    // ambil data siswa atau santri ke BE
    // dummy response

    console.log(qrDecoded);
    isDataReceived = true;
    dataBuyer = {
      name: "Ananda Navisa",
      nis: "123456789",
    };
  };

  const sendPurchase = async () => {
    if (!confirm("Konfirmasi?")) {
      return;
    }
    // kirim pembelian ke server
    isDataReceived = false;
    dataBuyer = {
      name: "",
      nis: "",
    };
  };

  const backToScan = () => {
    isDataReceived = false;
    dataBuyer = {
      name: "",
      nis: "",
    };
  };

  const stopScanning = () => {
    html5QrCode.stop().then(() => {
      console.log("The qr code scanning is stopped.");
    });
  };

  const startScanning = () => {
    html5QrCode
      .start(
        cameraId,
        {
          fps: 10, // Optional, frame per seconds for qr code scanning
          qrbox: { width: 250, height: 250 }, // Optional, if you want bounded box UI
        },
        (decodedText: string, decodedResult: QrcodeResult) => {
          stopScanning();

          getBuyerData(decodedText);
        },
        (err: any) => {
          console.log(err);
        },
      )
      .catch((err: any) => {
        // Start failed, handle it.
        console.log(err);
      });
  };
</script>

<section class="col border" style="height: 90vh;">
  {#if !isDataReceived}
    <div id="reader"></div>
    <div class="d-flex justify-content-center p-2 gap-3">
      <Button onclick={startScanning}>Scan</Button>
      <Button onclick={stopScanning}>Stop</Button>
    </div>
  {:else}
    <div class="d-flex justify-content-center p-2 gap-3">
      <div class="text-center">
        <p>Nama: {dataBuyer.name}</p>
        <p>Nis: {dataBuyer.nis}</p>
        <Button onclick={backToScan}>Back</Button>
        <Button onclick={sendPurchase}>Konfirmasi</Button>
      </div>
    </div>
  {/if}
</section>
