package markov3

import (
	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"
	"testing"
)

const lorem1 = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam gravida lobortis augue eu commodo. Aliquam pharetra mi eget ligula vehicula posuere. Donec sed augue viverra, eleifend ante ut, luctus arcu. Aliquam sagittis turpis eget tempus suscipit. Sed feugiat eros ac arcu lacinia, non tempus est tempor. Vestibulum orci nisl, cursus non interdum vitae, egestas id quam. Vestibulum sapien felis, rutrum a aliquam id, facilisis a turpis. Quisque pulvinar odio mi, et aliquet mi dignissim ut.`
const lorem2 = `In in nibh sodales lectus elementum ultricies id eget dolor. Aliquam efficitur massa ac ipsum aliquam vehicula. Aenean vitae nulla in ligula molestie posuere eu vitae turpis. Etiam a aliquam leo, at laoreet lacus. Nulla convallis dictum metus a lacinia. Nunc feugiat lectus nec lectus venenatis, sed ultrices felis ullamcorper. Nulla facilisis diam at tempus facilisis. Suspendisse rutrum tellus id massa eleifend ultricies. Proin laoreet, metus at tristique laoreet, ipsum elit ultrices neque, quis accumsan sem lacus at lacus. Cras elit leo, pharetra et est eget, pulvinar pellentesque ipsum. Curabitur blandit suscipit metus, eu bibendum neque laoreet ut. Phasellus hendrerit laoreet diam sit amet dignissim. Vestibulum tristique tempus nibh. Donec facilisis odio et pellentesque eleifend.`
const lorem3 = `Nulla facilisi. Suspendisse quis interdum libero. Duis non rhoncus lectus. Maecenas eleifend tortor ut nisl consequat ornare. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam et congue augue, in tempus libero. Praesent vel fringilla metus, eget ullamcorper justo. Morbi faucibus eleifend egestas. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla aliquet purus est, eget vestibulum nibh varius id. Mauris luctus sapien et augue pretium, in congue mi malesuada. Nullam vulputate lorem eu eros elementum semper. Aenean at nunc augue. Ut et urna non purus ultrices volutpat.`
const lorem4 = `Maecenas in efficitur lacus, eu facilisis mi. Nunc id purus faucibus, ullamcorper purus eu, efficitur mi. Sed sed posuere est. In pulvinar pellentesque nisl, vel interdum urna tempor eu. Mauris commodo, nunc at eleifend pulvinar, elit quam venenatis felis, sed venenatis erat est quis tellus. Praesent vitae mi a lectus ornare facilisis vestibulum et lacus. Nullam lacinia bibendum fringilla. In quis dolor a nisi vestibulum auctor. Donec efficitur felis ac sollicitudin laoreet. Phasellus molestie dictum consequat. Interdum et malesuada fames ac ante ipsum primis in faucibus. Aliquam lacinia dolor vel euismod suscipit. Proin nunc est, efficitur et ante sit amet, mattis lacinia nulla.`
const lorem5 = `Pellentesque tempor sem ut est tempor, at dictum magna cursus. Nunc nec interdum quam. Etiam nec vestibulum velit, quis facilisis velit. Aliquam pretium turpis at urna ultrices pellentesque. Sed quis rhoncus tellus, quis auctor nibh. Proin interdum erat sit amet augue vestibulum, sodales gravida neque placerat. Proin risus dui, commodo ut placerat ultrices, ornare non sem. In hac habitasse platea dictumst. Cras ut odio mi. Integer vitae lorem bibendum, blandit metus elementum, vestibulum sapien.`

func BenchmarkLearning(b *testing.B) {
	m := MakeMarkov("")
	for i := 0; i < b.N; i++ {
		m.LearnString("Hello, World!", cleaner.VariantPlain)
	}
}
