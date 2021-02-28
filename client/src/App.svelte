<script>
  import ProviderSearch from "./ProviderSearch.svelte";
</script>

<div class="container mx-auto pb-8">
  <header class="my-8 text-center">
    <h1 class="text-4xl font-serif">What Did Medicare Pay?</h1>
    <p class="my-2 md:w-1/2 max-w-lg mx-auto">
      Figure out how much a doctor was paid by Medicare for a certain procedure.
      Data pulled from the CMS 2017 public use file.
      <!-- Many providers (e.g. doctors) bill much more for a procedure than they expect to get paid.
      Insurers simply pay a pre-negotiated maximum billing amount, but uninsured patients may
      feel compelled to pay the full amount.  -->
    </p>
  </header>
  <ProviderSearch />
</div>
