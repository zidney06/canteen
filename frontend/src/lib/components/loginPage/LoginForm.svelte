<script lang="ts">
  import type { ErrorsType } from "$lib/types/types";
  import {
    Form,
    FormGroup,
    Input,
    FormFeedback,
    Button,
  } from "@sveltestrap/sveltestrap";

  interface Props {
    submit: () => void;
    selectedKasir: number | null;
    errors: ErrorsType;
    password: string;
    isSubmited: boolean;
    daftarKasir: number[];
  }

  let {
    submit,
    selectedKasir = $bindable(),
    errors,
    password = $bindable(),
    isSubmited,
    daftarKasir,
  }: Props = $props();
</script>

<Form onsubmit={submit}>
  <!-- RUBAH BAGIAN INI AGAR CLIENT HARUS MEMASUKAN CILENT KEY DAN CLIENT SECRET SECARA MANUAL -->
  <FormGroup>
    <Input
      type="select"
      bind:value={selectedKasir}
      required
      invalid={(errors.selectedKasir?.length ?? 0) > 0 && isSubmited}
    >
      {#each daftarKasir as option}
        <option value={option}>Kasir {option}</option>
      {/each}
    </Input>
    {#if (errors.selectedKasir?.length ?? 0) > 0}
      <FormFeedback>{errors.selectedKasir}</FormFeedback>
    {/if}
  </FormGroup>
  <FormGroup label="Password" floating>
    <Input
      placeholder="Password"
      type="password"
      bind:value={password}
      required
      invalid={(errors.password?.length ?? 0) > 0 && isSubmited}
    />
    {#if (errors.password?.length ?? 0) > 0}
      <FormFeedback>{errors.password}</FormFeedback>
    {/if}
  </FormGroup>
  <Button type="submit">Login</Button>
</Form>
