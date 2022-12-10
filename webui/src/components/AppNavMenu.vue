<template>
	<li class="nav-item">
		<RouterLink :to="to" class="nav-link">
			<div data-bs-target="#sidebarMenu.show, #sidebarMenu.collapsing" data-bs-toggle="collapse">
				<svg class="feather">
					<use :href="`${baseUrl}feather-sprite-v4.29.0.svg#${icon}`"/>
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
	data() {
		return {baseUrl: import.meta.env.BASE_URL};
	},
}
</script>

<style scoped>

</style>
