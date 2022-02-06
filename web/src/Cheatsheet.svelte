<script lang="ts">
  import SampleRow from "./SampleRow.svelte";

  export let samples;
</script>

<div class="lg:grid grid-cols-2 gap-4 p-4 sm:hidden">
  <table class="table-auto">
    <tbody>
      {#each samples.Left as sample}
        <SampleRow {sample} on:sample />
      {/each}
    </tbody>
  </table>

  <table class="table-auto">
    <tbody>
      {#each samples.Right as sample}
        <SampleRow {sample} on:sample />
      {/each}
    </tbody>
  </table>
</div>

<div class="gap-4 p-4 sm:block lg:hidden">
  <table class="table-auto w-full">
    <tbody>
      {#each samples.Left as sample}
        <SampleRow {sample} on:sample />
      {/each}
      {#each samples.Right as sample}
        <SampleRow {sample} on:sample />
      {/each}
    </tbody>
  </table>
</div>
