<script lang="ts">
  import favicon from "$lib/assets/favicon.svg";
  import {
    Collapse,
    Nav,
    Navbar,
    NavbarBrand,
    NavbarToggler,
    NavItem,
    NavLink,
    Styles,
  } from "@sveltestrap/sveltestrap";

  let { children } = $props();
  let isOpen = $state(false);

  const handleUpdate = (event: any) => {
    isOpen = event.detail.isOpen;
  };
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

<Styles />

<Navbar color="secondary" theme="dark">
  <NavbarBrand href="/">Canteen</NavbarBrand>
  <NavbarToggler onclick={() => (isOpen = !isOpen)} />
  <Collapse {isOpen} navbar expand="md" on:update={handleUpdate}>
    <Nav class="ms-auto" navbar>
      <NavItem>
        <NavLink href="/app/dashboard">Dashboard</NavLink>
      </NavItem>
      <NavItem>
        <NavLink href="/app/product">Produk</NavLink>
      </NavItem>
    </Nav>
  </Collapse>
</Navbar>
{@render children()}
