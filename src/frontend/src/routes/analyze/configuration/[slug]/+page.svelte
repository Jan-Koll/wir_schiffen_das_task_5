<script lang="ts">
	import type { PageData } from './$types';
	import type { configuration, engine_control_analyzer_response, cooling_exhaust_analyzer_response, mounting_analyzer_response, propulsion_analyzer_response, supply_analyzer_response } from '$lib/types';

	import { page } from '$app/stores';
	export let data: {
		configuration: configuration;
		cooling_exhaust_analyzer: cooling_exhaust_analyzer_response[];
		engine_control_analyzer: engine_control_analyzer_response[];
		mounting_analyzer: mounting_analyzer_response[];
		propulsion_analyzer: propulsion_analyzer_response[];
		supply_analyzer: supply_analyzer_response[];
	};
	import ChevronLeft from 'lucide-svelte/icons/chevron-left';
	import LoaderCircle from 'lucide-svelte/icons/loader-circle';
	import Play from 'lucide-svelte/icons/play';
	import CircleCheck from 'lucide-svelte/icons/circle-check';
	import CircleX from 'lucide-svelte/icons/circle-x';
	import Circle from 'lucide-svelte/icons/circle';
	import CircleEllpsis from 'lucide-svelte/icons/circle-ellipsis';

	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { goto } from '$app/navigation';
	import { writable, derived } from 'svelte/store';
	import { toast } from 'svelte-sonner';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { messages } from '$lib/state/sse';
	import StatusIndicator from '$lib/components/StatusIndicator.svelte';
	let awaitingResponse = false;
	let latestAnalysisId = data.engine_control_analyzer.length > 1 ? data.engine_control_analyzer[data.engine_control_analyzer.length - 1].id : 1;

	let latestCoolingExhaustAnalysis = data.cooling_exhaust_analyzer.find(
		(result) => result.id === latestAnalysisId
	);
	let latestEngineControlAnalysis = data.engine_control_analyzer.find(
		(result) => result.id === latestAnalysisId
	);
	let latestMountingAnalysis = data.mounting_analyzer.find(
		(result) => result.id === latestAnalysisId
	);
	let latestPropulsionAnalysis = data.propulsion_analyzer.find(
		(result) => result.id === latestAnalysisId
	);
	let latestSupplyAnalysis = data.supply_analyzer.find(
		(result) => result.id === latestAnalysisId
	);
	/*$: if (latestAnalysisId) {
		latestEngineControlAnalysis = data.engine_control_analyzer.find(
			(result) => result.id === latestAnalysisId
		);
	}*/
	console.log('loaded data', data);

	function updateAnalysisState(parsedMessage: any) {
		//todo: irgendwie abstrahieren, dass es weniger repetativer code ist
		switch (parsedMessage.sender) {
			case 'cooling_exhaust_analyzer': {
				const idx = data.cooling_exhaust_analyzer.findLastIndex(
					(analysis_result) => analysis_result.id === parsedMessage.id
				);
				if (idx == -1) {
					//insert
					data.cooling_exhaust_analyzer.push(parsedMessage);
					data.cooling_exhaust_analyzer = data.cooling_exhaust_analyzer;
					if (parsedMessage.id >= latestAnalysisId) {
						latestCoolingExhaustAnalysis = parsedMessage
						latestAnalysisId = parsedMessage.id;
					}
				} else {
					//update
					data.cooling_exhaust_analyzer[idx] = parsedMessage;
					latestCoolingExhaustAnalysis = parsedMessage
				}
				break;
			}
			case 'engine_control_analyzer': {
				const idx = data.engine_control_analyzer.findLastIndex(
					(analysis_result) => analysis_result.id === parsedMessage.id
				);
				if (idx == -1) {
					//insert
					data.engine_control_analyzer.push(parsedMessage);
					data.engine_control_analyzer = data.engine_control_analyzer;
					if (parsedMessage.id >= latestAnalysisId) {
						latestEngineControlAnalysis = parsedMessage
						latestAnalysisId = parsedMessage.id;
					}
				} else {
					//update
					data.engine_control_analyzer[idx] = parsedMessage;
					latestEngineControlAnalysis = parsedMessage
				}
				break;
			}
			case 'mounting_analyzer': {
				const idx = data.mounting_analyzer.findLastIndex(
					(analysis_result) => analysis_result.id === parsedMessage.id
				);
				if (idx == -1) {
					//insert
					data.mounting_analyzer.push(parsedMessage);
					data.mounting_analyzer = data.mounting_analyzer;
					if (parsedMessage.id >= latestAnalysisId) {
						latestMountingAnalysis = parsedMessage
						latestAnalysisId = parsedMessage.id;
					}
				} else {
					//update
					data.mounting_analyzer[idx] = parsedMessage;
					latestMountingAnalysis = parsedMessage
				}
				break;
			}

			case 'propulsion_analyzer': {
				const idx = data.propulsion_analyzer.findLastIndex(
					(analysis_result) => analysis_result.id === parsedMessage.id
				);
				if (idx == -1) {
					//insert
					data.propulsion_analyzer.push(parsedMessage);
					data.propulsion_analyzer = data.propulsion_analyzer;
					if (parsedMessage.id >= latestAnalysisId) {
						latestPropulsionAnalysis = parsedMessage
						latestAnalysisId = parsedMessage.id;
					}
				} else {
					//update
					data.propulsion_analyzer[idx] = parsedMessage;
					latestPropulsionAnalysis = parsedMessage
				}
				break;
			}
			case 'supply_analyzer': {
				const idx = data.supply_analyzer.findLastIndex(
					(analysis_result) => analysis_result.id === parsedMessage.id
				);
				if (idx == -1) {
					//insert
					data.supply_analyzer.push(parsedMessage);
					data.supply_analyzer = data.supply_analyzer;
					if (parsedMessage.id > latestAnalysisId) {
						latestSupplyAnalysis = parsedMessage
						latestAnalysisId = parsedMessage.id;
					}
				} else {
					//update
					data.supply_analyzer[idx] = parsedMessage;
					latestSupplyAnalysis = parsedMessage
				}
				break;
			}
		}
	}
	$: if ($messages) {
		let parsedMessage;
		console.log('message: ', $messages);
		try {
			parsedMessage = JSON.parse($messages);
			parsedMessage = parsedMessage.message;
			if (!parsedMessage.sender) {
				toast($messages);
			} else if ('configuration_id' in parsedMessage)
				if (parsedMessage.configuration_id === data.configuration.id) {
					//own configuration
					updateAnalysisState(parsedMessage);
				} else {
					// other configurations
					toast($messages);
				}
		} catch (e) {
			console.log('error', e);
			toast($messages);
		}
	}
	async function startAnalyzing() {
		let formData = { configuration_id: data.configuration.id };
		const body = JSON.stringify(formData);
		awaitingResponse = true;
		fetch($page.url.pathname, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body
		})
			.then((response) => {
				console.log(response);
				if (response.status >= 300) {
					toast.error('Error ' + response.status, {
						description: 'Failed starting the analysis'
					});
					console.log(`Error ${response.status}: Failed starting the analysis.`);
					return false;
				}
				return response.json();
			})
			.then((analysesData) => {
				awaitingResponse = false;
				if (analysesData) {
					toast.info('Analysis added to queue');
					console.log(
						'Analysis successfully added to queue:',
						analysesData,
						data.engine_control_analyzer
					);
				}
			})
			.catch((error) => {
				awaitingResponse = false;
				console.error('Error starting analysis:', error);
				toast.error('Error', { description: 'Could not start analysis: ' + error });
			});
	}
</script>

<div class="mx-auto grid max-w-[59rem] flex-1 auto-rows-max gap-4">
	<div class="flex items-center gap-4">
		<Button variant="outline" size="icon" class="h-7 w-7" href="/analyze">
			<ChevronLeft class="h-4 w-4" />
			<span class="sr-only">Back</span>
		</Button>
		<h1 class="flex-1 shrink-0 whitespace-nowrap text-xl font-semibold tracking-tight sm:grow-0">
			View Configuration #{data.configuration.id}
		</h1>
		<div class="hidden items-center gap-2 md:ml-auto md:flex">
			<Button
				type="reset"
				variant="outline"
				size="sm"
				href={`/analyze/edit/${data.configuration.id}`}>Edit</Button
			>
			<Button size="sm" disabled={awaitingResponse} on:click={startAnalyzing}>
				{#if !awaitingResponse}
					<Play class="mr-2 h-4 w-4"></Play>
					Start Analysis
				{:else}
					<LoaderCircle class="mr-2 h-4 w-4 animate-spin" />
					Please wait
				{/if}
			</Button>
		</div>
		<!-- todo: when viewing, show badge with id: -->
	</div>
	<div class="grid gap-4 md:grid-cols-[1fr_250px] lg:grid-cols-3 lg:gap-8">
		<div class="grid-auto-rows-max grid items-start gap-4 lg:col-span-2 lg:gap-8">
			<Card.Root>
				<Card.Header>
					<Card.Title>General</Card.Title>
					<Card.Description
						>General information about the configuration and the used engine</Card.Description
					>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-6">
						<div class="grid gap-3" class:opacity-60={data.configuration.description === ''}>
							<Label>Description <span class="text-xs">(required)</span></Label>
							{#if data.configuration.description !== ''}
								{data.configuration.description}
							{:else}
								<span class="text-red-500">No description</span>
							{/if}
						</div>
						<!--<div class="flex flex-col gap-6 justify-normal space-y-6 sm:flex-row sm:space-y-0">-->
						<div class="grid gap-8 sm:grid-cols-2">
							<div class="w grid gap-3" class:opacity-60={data.configuration.engine === ''}>
								<Label>Engine <span class="text-xs">(required)</span></Label>
								{#if data.configuration.engine !== ''}
									{data.configuration.engine}
								{:else}
									<span class="text-red-500">No engine</span>
								{/if}
							</div>
							<div class="grid gap-3" class:opacity-60={data.configuration.gearbox_type === ''}>
								<Label for="gearboxType">Gearbox type <span class="text-xs">(required)</span></Label
								>
								{#if data.configuration.gearbox_type !== ''}
									{data.configuration.gearbox_type}
								{:else}
									<span class="text-red-500">No gearbox type</span>
								{/if}
							</div>
						</div>
					</div>
				</Card.Content>
			</Card.Root>
			<Card.Root>
				<Card.Header>
					<Card.Title>Engine Control</Card.Title>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-6">
						<div class="grid gap-3" class:opacity-60={!data.configuration.engine_management_system}>
							<Label>Engine Management System</Label>
							<div class="flex items-center space-x-2">
								<Checkbox
									id="engineManagementSystem"
									bind:checked={data.configuration.engine_management_system}
									aria-labelledby="engineManagementSystemLabel"
									disabled
								/>
								<Label
									id="engineManagementSystemLabel"
									for="engineManagementSystem"
									class="my-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-100"
								>
									Compliance with Classification Society Regulations
								</Label>
							</div>
						</div>
						<div
							class="grid gap-3"
							class:opacity-60={data.configuration.monitoring_control_system === ''}
						>
							<Label for="monitoringControlSystem">Monitoring/Control System</Label>
							{#if data.configuration.monitoring_control_system !== ''}
								{data.configuration.monitoring_control_system}
							{:else}
								None
							{/if}
						</div>
					</div>
				</Card.Content>
			</Card.Root>
			<Card.Root>
				<Card.Header>
					<Card.Title>Propulsion</Card.Title>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-6">
						<div class="grid gap-3" class:opacity-60={!data.configuration.starting_system}>
							<Label>Starting System</Label>
							<div class="flex items-center space-x-2">
								<Checkbox
									disabled
									id="startingSystem"
									bind:checked={data.configuration.starting_system}
									aria-labelledby="startingSystemLabel"
								/>
								<Label
									id="startingSystemLabel"
									for="startingSystem"
									class="my-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-100"
									>Air starter
								</Label>
							</div>
						</div>
						<div class="grid gap-3" class:opacity-60={!data.configuration.power_transmission}>
							<Label>Power Transmission</Label>
							<div class="flex items-center space-x-2">
								<Checkbox
									id="powerTransmission"
									disabled
									bind:checked={data.configuration.power_transmission}
									aria-labelledby="powerTransmissionLabel"
								/>
								<Label
									id="powerTransmissionLabel"
									for="powerTransmission"
									class="my-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-100"
									>Torsionally resilient coupling</Label
								>
							</div>
						</div>
						<div
							class="grid gap-3"
							class:opacity-60={!data.configuration.auxiliary_pto ||
								(data.configuration.auxiliary_pto && data.configuration.auxiliary_pto.length === 0)}
						>
							<Label for="auxiliaryPto">Auxiliary PTO</Label>
							{#if data.configuration.auxiliary_pto && data.configuration.auxiliary_pto.length > 0}
								<ul class="ml-5 list-disc">
									{#each data.configuration.oil_system as item}
										<li>{item}</li>
									{/each}
								</ul>
							{:else}
								None
							{/if}
						</div>
					</div>
				</Card.Content>
			</Card.Root>
			<Card.Root>
				<Card.Header>
					<Card.Title>Supply</Card.Title>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-6">
						<div
							class="grid gap-3"
							class:opacity-60={!data.configuration.oil_system ||
								(data.configuration.oil_system && data.configuration.oil_system.length === 0)}
						>
							<Label for="oilSystems">Oil system</Label>
							{#if data.configuration.oil_system && data.configuration.oil_system.length > 0}
								<ul class="ml-5 list-disc">
									{#each data.configuration.oil_system as item}
										<li>{item}</li>
									{/each}
								</ul>
							{:else}
								None
							{/if}
						</div>
						<div
							class="grid gap-3"
							class:opacity-60={!data.configuration.fuel_system ||
								(data.configuration.fuel_system && data.configuration.fuel_system.length === 0)}
						>
							<Label for="fuelSystems">Fuel system</Label>
							{#if data.configuration.fuel_system && data.configuration.fuel_system.length > 0}
								<ul class="ml-5 list-disc">
									{#each data.configuration.fuel_system as item}
										<li>{item}</li>
									{/each}
								</ul>
							{:else}
								None
							{/if}
						</div>
					</div>
				</Card.Content>
			</Card.Root>
			<Card.Root>
				<Card.Header>
					<Card.Title>Cooling & Exhaust</Card.Title>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-6">
						<div class="grid gap-3" class:opacity-60={data.configuration.cooling_system === ''}>
							<Label for="coolingSystem">Cooling system</Label>
							{#if data.configuration.cooling_system !== ''}
								{data.configuration.cooling_system}
							{:else}
								None
							{/if}
						</div>
						<div class="grid gap-3" class:opacity-60={!data.configuration.exhaust_system}>
							<Label>Exhaust system</Label>
							<div class="flex items-center space-x-2">
								<Checkbox
									id="exhaustSystem"
									disabled
									bind:checked={data.configuration.exhaust_system}
									aria-labelledby="exhaustSystemLabel"
								/>
								<Label
									id="exhaustSystemLabel"
									for="exhaustSystem"
									class="my-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-100"
									>90° Exhaust bellows discharge rotatable</Label
								>
							</div>
						</div>
					</div>
				</Card.Content>
			</Card.Root>
			<Card.Root>
				<Card.Header>
					<Card.Title>Mounting & Gearbox</Card.Title>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-6">
						<div class="grid gap-3" class:opacity-60={!data.configuration.mounting_system}>
							<Label>Mounting system</Label>
							<div class="flex items-center space-x-2">
								<Checkbox
									disabled
									id="mountingSystem"
									bind:checked={data.configuration.mounting_system}
									aria-labelledby="mountingSystemLabel"
								/>
								<Label
									id="mountingSystemLabel"
									for="mountingSystem"
									class="my-2 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-100"
									>Resilient mounts at driving end</Label
								>
							</div>
						</div>
						<div
							class="grid gap-3"
							class:opacity-60={!data.configuration.gearbox_option ||
								(data.configuration.gearbox_option &&
									data.configuration.gearbox_option.length === 0)}
						>
							<Label for="gearboxOptions">Gearbox options</Label>
							{#if data.configuration.gearbox_option && data.configuration.gearbox_option.length > 0}
								<ul class="ml-5 list-disc">
									{#each data.configuration.gearbox_option as item}
										<li>{item}</li>
									{/each}
								</ul>
							{:else}
								None
							{/if}
						</div>
					</div>
				</Card.Content>
			</Card.Root>
		</div>

		<div class="grid auto-rows-max items-start gap-4 lg:gap-8">
			<Card.Root>
				<Card.Header>
					<Card.Title>Status #{latestAnalysisId}</Card.Title>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-6">
						<div class="grid gap-3">
							{#if data.engine_control_analyzer.length > 0}
								<div class="flex gap-3">
									<StatusIndicator
										status={latestEngineControlAnalysis?.status}
										success={latestEngineControlAnalysis?.success}
										name="Engine Control Analyzer #{latestAnalysisId}"
										side="left"
									></StatusIndicator>
									Engine Control Analyzer
								</div>
							{/if}
							{#if data.cooling_exhaust_analyzer.length > 0}
								<div class="flex gap-3">
									<StatusIndicator
										status={latestCoolingExhaustAnalysis?.status}
										success={latestCoolingExhaustAnalysis?.success}
										name="Cooling Exhaust Analyzer #{latestAnalysisId}"
										side="left"
									></StatusIndicator>
									Cooling Exhaust Analyzer
								</div>
							{/if}

							{#if data.mounting_analyzer.length > 0}
								<div class="flex gap-3">
									<StatusIndicator
										status={latestMountingAnalysis?.status}
										success={latestMountingAnalysis?.success}
										name="Mounting Analyzer #{latestAnalysisId}"
										side="left"
									></StatusIndicator>
									Mounting Analyzer
								</div>
							{/if}

							{#if data.propulsion_analyzer.length > 0}
								<div class="flex gap-3">
									<StatusIndicator
										status={latestPropulsionAnalysis?.status}
										success={latestPropulsionAnalysis?.success}
										name="Propulsion Analyzer #{latestAnalysisId}"
										side="left"
									></StatusIndicator>
									Propulsion Analyzer
								</div>
							{/if}

							{#if data.supply_analyzer.length > 0}
								<div class="flex gap-3">
									<StatusIndicator
										status={latestSupplyAnalysis?.status}
										success={latestSupplyAnalysis?.success}
										name="Supply Analyzer #{latestAnalysisId}"
										side="left"
									></StatusIndicator>
									Supply Analyzer
								</div>
							{/if}

							<!--<StatusIndicator status="running" success={false} name="Test" side="left"></StatusIndicator>
							<CircleCheck class="rounded-full text-emerald-400"></CircleCheck>
							<CircleX class="rounded-full text-red-400"></CircleX>
							<CircleEllpsis class="rounded-full text-sky-500"></CircleEllpsis>
							<Circle class="rounded-full text-gray-400"></Circle>-->
						</div>
					</div>
				</Card.Content>
			</Card.Root>
			<Card.Root class="overflow-hidden">
				<Card.Header>
					<Card.Title>Result history</Card.Title>
					<Card.Description>Results of past analyses of this configuration</Card.Description>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-2">
						<div
							class="light:border-emerald-400 flex items-center gap-x-1 rounded border p-1 dark:border-emerald-800"
						>
							<span class="text-muted-foreground mx-2">#1</span>
							<Tooltip.Root>
								<Tooltip.Trigger>
									<CircleCheck class="h-5 w-5 rounded-full text-emerald-400"></CircleCheck>
								</Tooltip.Trigger>
								<Tooltip.Content>
									<p>Cooling-Exhaust-Analyzer</p>
								</Tooltip.Content>
							</Tooltip.Root>
							<CircleX class="h-5 w-5 rounded-full text-red-400"></CircleX>
							<CircleEllpsis class="h-5 w-5 rounded-full text-sky-500"></CircleEllpsis>
							<Circle class="h-5 w-5 rounded-full text-gray-400"></Circle>
							<Circle class="h-5 w-5 rounded-full text-gray-400"></Circle>
						</div>
					</div>
				</Card.Content>
			</Card.Root>
		</div>
	</div>
	<!--<div class="flex items-center justify-center gap-2 md:justify-end">
				<Button type="reset" variant="outline" href={`/analyze/edit/${data.configuration.id}`}>Edit</Button>
				<Button disabled={awaitingResponse} on:click={startAnalyzing}>
					{#if !awaitingResponse}
						Analyze
					{:else}
						<LoaderCircle class="mr-2 h-4 w-4 animate-spin" />
						Please wait
					{/if}
				</Button>
			</div>-->
</div>
