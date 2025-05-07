<script lang="ts">
	import Icon from "@iconify/svelte";
	import { themeChange } from "theme-change";
	$effect(() => {
		themeChange(false);
	});
</script>

<div>
	<label class="swap swap-rotate">
		<input type="checkbox" class="hidden" data-toggle-theme="dark,nord" />
		<Icon icon="ph:sun" class="icon24 swap-off" />
		<Icon icon="ph:moon" class="icon24  swap-on" />
	</label>
</div>
