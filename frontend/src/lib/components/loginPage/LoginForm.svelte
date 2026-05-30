<script lang="ts">
  import { Form, FormGroup, Input, Button } from "@sveltestrap/sveltestrap";

  interface Props {
    submit: () => void;
    clientKey: string;
    clientSecret: string;
  }

  let {
    submit,
    clientKey = $bindable(),
    clientSecret = $bindable(),
  }: Props = $props();
</script>

<Form onsubmit={submit}>
  <FormGroup label="clientKey" floating>
    <Input
      type="text"
      placeholder="clientKey"
      bind:value={clientKey}
      required
    />
  </FormGroup>
  <FormGroup label="clientSecret" floating>
    <Input
      placeholder="clientSecret"
      type="password"
      bind:value={clientSecret}
      required
    />
  </FormGroup>
  <Button type="submit">Login</Button>
</Form>
