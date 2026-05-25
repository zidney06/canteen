<script lang="ts">
  import favicon from "$lib/assets/favicon.svg";
  import { Styles } from "@sveltestrap/sveltestrap";

  let { children } = $props();
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

<Styles />
{@render children()}
