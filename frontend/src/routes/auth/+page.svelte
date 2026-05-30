<script lang="ts">
  import { goto } from "$app/navigation";
  import LoginForm from "$lib/components/loginPage/LoginForm.svelte";
  import { clientLogin } from "$lib/repo/repository";
  import type { ErrorsType } from "$lib/types/types";
  import {
    Card,
    CardBody,
    CardHeader,
    CardTitle,
    Container,
  } from "@sveltestrap/sveltestrap";
  import { onMount } from "svelte";
  import { z } from "zod";

  /*
  ck: CK_b4a62303-dd5c-4e21-a0d8-86c486da8f04
  cs: 4573e53b01cd1c921b93ad1d85502cda8d3d63be7fb3bb57baeba679187651d0
  */

  const reqBodySchema = z.object({
    clientSecret: z.string(),
    clientKey: z.string().min(39, "Client Key minimal 39 karakter"),
  });

  let clientKey: string = $state("");
  let clientSecret: string = $state("");

  // ?
  let errors: ErrorsType = $state({
    clientSecret: [],
    clientKey: [],
  });

  onMount(() => {
    // ambil daftar kasir di sini
  });

  // client secret & key disimpan di cache ketika membuka one time url
  const login = (clientecret: string, clientKey: string) => {
    if (clientKey.trim() == "" && clientecret.trim() == "") {
      alert("Client Secret dan Client Key tidak boleh kosong");
      return;
    }
  };

  const submit = async () => {
    errors = {
      clientSecret: [],
      clientKey: [],
    };

    const result = reqBodySchema.safeParse({
      clientSecret,
      clientKey,
    });

    console.log(result, clientSecret, clientKey);

    if (result.success) {
      // request ke backend. saat berhasil login redirect ke dashboard
      const res = await clientLogin(clientKey, clientSecret);

      if (!res.isSuccess) {
        alert(res.message);
        return;
      }

      console.log(res);

      // simpan jwt
      sessionStorage.setItem("canteen-jwt", res.jwtToken);

      clientSecret = "";
      clientKey = "";

      goto("/app/dashboard");
    } else {
      errors = z.flattenError(result.error).fieldErrors;
      alert(`
        Error:
          field client key: ${errors.clientKey ?? "Tidak ada error"}
          field client secret: ${errors.clientSecret ?? "Tidak ada error"}
        `);
    }
  };
</script>

<Container
  class="d-flex justify-content-center align-items-center"
  style="height: 100vh;"
>
  <Card style="width: 450px;">
    <CardHeader>
      <CardTitle class="text-center">Login Kasir</CardTitle>
    </CardHeader>
    <CardBody>
      <LoginForm bind:clientKey bind:clientSecret {submit} />
    </CardBody>
  </Card>
</Container>
