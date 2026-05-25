<script lang="ts">
  import LoginForm from "$lib/components/loginPage/LoginForm.svelte";
  import type { ErrorsType } from "$lib/types/types";
  import {
    Card,
    CardBody,
    CardHeader,
    CardTitle,
    Container,
    FormGroup,
    Input,
    Form,
    Button,
    FormFeedback,
  } from "@sveltestrap/sveltestrap";
  import { onMount } from "svelte";
  import { z } from "zod";

  const reqBodySchema = z.object({
    password: z.string().min(8, { message: "Minimal 8 karakter" }),
    selectedKasir: z.number().refine((val) => daftarKasir.includes(val), {
      message: "Value kasir harus salah satu dari [1, 2, 3]",
    }),
  });

  let selectedKasir: number | null = $state(null);
  const daftarKasir: number[] = [1, 2, 3]; // ini ambil dari server aja
  let password: string = $state("");
  let isSubmited: boolean = $state(false);
  let errors: ErrorsType = $state({
    password: [],
    selectedKasir: [],
  });

  onMount(() => {
    // ambil daftar kasir di sini
  });

  const submit = () => {
    isSubmited = true;
    errors = {
      password: [],
      selectedKasir: [],
    };

    const result = reqBodySchema.safeParse({
      password,
      selectedKasir,
    });

    console.log(result, password);

    if (result.success) {
      // request ke backend. saat berhasil login redirect ke dashboard
      console.log("Login berhasil");

      password = "";
      selectedKasir = null;
    } else {
      errors = z.flattenError(result.error).fieldErrors;
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
      <LoginForm
        bind:selectedKasir
        {daftarKasir}
        {isSubmited}
        {errors}
        bind:password
        {submit}
      />
    </CardBody>
  </Card>
</Container>
