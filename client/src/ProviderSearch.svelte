<script>
  import { titleCase } from "./utils/text";

  let query: string;
  let pendingSearchResults: Promise<{ Results: Array<any | null> }>;
  async function fetchPUPDResults(query: string) {
    const resp = await fetch(`/api/search?query=${encodeURIComponent(query)}`);
    return resp.json();
  }
  function handleSearch() {
    pendingSearchResults = fetchPUPDResults(query);
  }
  function handleDemoSearch(demoQuery: string) {
    return () => {
      query = demoQuery;
      handleSearch();
    };
  }
</script>

<form
  on:submit|preventDefault={handleSearch}
  class="w-11/12 md:w-2/3 lg:w-1/2 my-8 mb-4 mx-auto flex flex-col"
>
  <div class="flex flex-row">
    <input
      bind:value={query}
      class="w-full mx-2 p-2 border border-transparent bg-green-50 rounded ring-2 ring-green-600 focus:outline-none"
    />
    <button
      class="mx-2 p-2 rounded text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-600 focus:ring-opacity-50"
    >
      Search
    </button>
  </div>
  <small class="mx-2 my-1 italic">
    try searching "<a
      href="#"
      class="underline"
      on:click|preventDefault={handleDemoSearch("robert anesthesia")}
    >
      robert anesthesia</a
    >" or "<a
      href="#"
      class="underline"
      on:click|preventDefault={handleDemoSearch("mri address:CA")}
    >
      mri address:CA</a
    >"
  </small>
</form>
<table class="w-11/12 mx-auto table-auto">
  <thead class="text-xs">
    <tr>
      <th class="text-green-600">Provider</th>
      <th class="text-green-600">Address</th>
      <th class="text-green-600">Procedure</th>
      <th class="text-green-600">Billed Cost (# of records)</th>
      <th class="text-green-600">Real Cost (& Copay)</th>
      <th class="text-green-600">Adjustment</th>
    </tr>
  </thead>
  <tbody class="text-xs">
    {#if pendingSearchResults}
      {#await pendingSearchResults}
        <pre>Loading...</pre>
      {:then searchResp}
        {#if searchResp.Results}
          {#each searchResp.Results as pupd (pupd.RowID)}
            <tr>
              <td class="p-2 border border-green-600 text-green-600">
                {titleCase(
                  pupd.NPPESProviderFirstName +
                    " " +
                    pupd.NPPESProviderLastOrgName
                )}
              </td>
              <td class="p-2 border border-green-600 text-green-600">
                {titleCase(
                  pupd.NPPESProviderStreet1 +
                    (pupd.NPPESProviderStreet2
                      ? " " + pupd.NPPESProviderStreet2
                      : "") +
                    ", " +
                    pupd.NPPESProviderCity
                )},
                {pupd.NPPESProviderState}
              </td>
              <td class="p-2 border border-green-600 text-green-600">
                {pupd.HCPCSCode}: {pupd.HCPCSDescription} ({pupd.PlaceOfService ==
                "F"
                  ? "facility"
                  : "office"})
              </td>
              <td class="p-2 border border-green-600 text-green-600">
                ${parseFloat(pupd.AverageSubmittedChrgAmt).toFixed(2)}<br />
                ({pupd.LineSrvcCnt})
              </td>
              <td class="p-2 border border-green-600 text-green-600">
                <!-- TODO: handle negative copay -->
                ${parseFloat(pupd.AverageMedicarePaymentAmt).toFixed(2)} + ${(
                  parseFloat(pupd.AverageMedicareAllowedAmt) -
                  parseFloat(pupd.AverageMedicarePaymentAmt)
                ).toFixed(2)}
              </td>
              <td
                class="p-2 text-center border border-green-600 text-green-600"
              >
                {(
                  (1 -
                    parseFloat(pupd.AverageMedicareAllowedAmt) /
                      parseFloat(pupd.AverageSubmittedChrgAmt)) *
                  -100
                ).toFixed(0)}%
              </td>
            </tr>
          {/each}
        {:else}
          <tr><td class="text-center" colspan="20">No results found.</td></tr>
        {/if}
      {/await}
    {/if}
  </tbody>
</table>
<!-- {#if pendingSearchResults}
  {#await pendingSearchResults then searchResp}
    <pre
      class="overflow-scroll w-11/12 mx-auto">{JSON.stringify(searchResp, null, 2)}</pre>
  {:catch error}
    <pre>{error}</pre>
  {/await}
{/if} -->
