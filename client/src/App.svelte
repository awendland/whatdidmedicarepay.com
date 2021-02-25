<script>
  let query = "";
  let queryDetails: Array<{ word: string; length: number }> = [];
  function handleSearch() {
    queryDetails = query
      .split(" ")
      .map((word) => ({ word, length: word.length }));
  }
</script>

<div class="container mx-auto">
  <header class="my-8 text-center">
    <h1 class="text-4xl font-serif">What Did Medicare Pay?</h1>
    <p class="my-2 md:w-1/2 max-w-lg mx-auto">
      Figure out how much a doctor was paid by Medicare for a certain procedure.
      Data pulled from the CMS 2017 public use file.
    </p>
  </header>
  <form
    on:submit|preventDefault={handleSearch}
    class="w-11/12 md:w-1/2 my-8 mx-auto flex flex-row"
  >
    <input
      bind:value={query}
      class="mx-2 w-full p-2 border border-transparent bg-green-50 rounded ring-2 ring-green-600 focus:outline-none"
    />
    <button
      class="mx-2 p-2 rounded text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-600 focus:ring-opacity-50"
    >
      Search
    </button>
  </form>
  <table class="w-11/12 md:w-2/3 mx-auto table-auto">
    <thead>
      <tr>
        <th class="text-green-600">Word</th>
        <th class="text-green-600">Length</th>
      </tr>
    </thead>
    <tbody>
      {#each queryDetails as detail}
        <tr>
          <td class="p-2 border border-green-600 text-green-600">
            {detail.word}
          </td>
          <td class="p-2 border border-green-600 text-green-600">
            {detail.length}
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
