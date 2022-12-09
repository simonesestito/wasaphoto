<template>
	<li class="nav-item">
		<RouterLink :to="to" class="nav-link">
			<div data-bs-toggle="collapse" data-bs-target="#sidebarMenu.show, #sidebarMenu.collapsing">
				<svg class="feather">
					<use :href="`/feather-sprite-v4.29.0.svg#${icon}`"/>
				</svg>
				{{ text }}
			</div>
		</RouterLink>
	</li>
</template>

<script>
export default {
	name: "AppNavMenu",
	props: ['to', 'text', 'icon'],
}
</script>

<style scoped>

</style>
