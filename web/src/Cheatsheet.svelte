<script lang="ts">
  import SampleRow from "./SampleRow.svelte";
  import { samplesLeft, samplesRight } from "./samples";
</script>

<div class="lg:grid grid-cols-2 gap-4 p-4 sm:hidden">
  <table class="table-auto">
    <tbody>
      {#each samplesLeft as sample}
        <SampleRow {sample} on:sample />
      {/each}
    </tbody>
  </table>

  <table class="table-auto">
    <tbody>
      {#each samplesRight as sample}
        <SampleRow {sample} on:sample />
      {/each}
    </tbody>
  </table>
</div>

<div class="gap-4 p-4 sm:block lg:hidden">
  <table class="table-auto w-full">
    <tbody>
      {#each samplesLeft as sample}
        <SampleRow {sample} on:sample />
      {/each}
      {#each samplesRight as sample}
        <SampleRow {sample} on:sample />
      {/each}
    </tbody>
  </table>
</div>
