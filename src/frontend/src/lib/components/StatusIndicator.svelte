<script lang="ts">
	import CircleCheck from 'lucide-svelte/icons/circle-check';
	import CircleX from 'lucide-svelte/icons/circle-x';
	import Circle from 'lucide-svelte/icons/circle';
	import CircleEllpsis from 'lucide-svelte/icons/circle-ellipsis';
	import CircleAlert from 'lucide-svelte/icons/circle-alert';
	import CircleSlash from 'lucide-svelte/icons/circle-slash';
	import * as Tooltip from '$lib/components/ui/tooltip';
	export let status: string|undefined;
	export let success: boolean=false;
	export let side: string = "side"
	export let name: string
	//export let small: boolean = false;
</script>

<Tooltip.Root openDelay={0}>
	<Tooltip.Trigger>
		{#if status == 'queued'}
			<Circle class="rounded-full text-gray-400"></Circle>
		{:else if status == 'running'}
			<CircleEllpsis class="animate-pulse rounded-full text-sky-500"></CircleEllpsis>
		{:else if status == 'canceled'}
			<CircleSlash class="rounded-full text-yellow-600"></CircleSlash>
		{:else if status == 'failed'}
			<CircleAlert class="rounded-full text-red-600"></CircleAlert>
		{:else if status == 'ready'}
			{#if success}
				<CircleCheck class="rounded-full text-emerald-400"></CircleCheck>
			{:else}
				<CircleX class="rounded-full text-red-400"></CircleX>
			{/if}
		{/if}
	</Tooltip.Trigger>
	<Tooltip.Content side={side}>
			<div class="font-bold mb-3">{name}</div>
			<div class="flex flex-col">
			<div><span class="font-bold">Status: </span>{status}</div>
			{#if status == 'ready'}
				<div><span class="font-bold">Success: </span>{success}</div>
			{/if}
		</div>
	</Tooltip.Content>
</Tooltip.Root>
