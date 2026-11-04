package search

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
)

func titled(id, title, cue string) lesson.Lesson {
	return lesson.Lesson{ID: id, Draft: lesson.Draft{Title: title, Cue: cue}}
}

func ids(lessons []lesson.Lesson) []string {
	out := make([]string, 0, len(lessons))
	for _, l := range lessons {
		out = append(out, l.ID)
	}

	return out
}

func TestSimilarRanksBySharedDistinctiveWords(t *testing.T) {
	draft := lesson.Draft{Title: "Run migrations before seeding the database", Cue: "make seed fails with relation users missing"}
	candidates := []lesson.Lesson{
		titled("a", "Seed the database after migrations", "make seed fails"),
		titled("b", "Install the toolchain with mise", "go is missing"),
		titled("c", "Migrations need a running database", "relation users does not exist"),
		titled("d", "Run the linter", "the build fails"),
	}

	assert.Equal(t, []string{"a", "c"}, ids(Similar(draft, candidates, 5)))
}

func TestSimilarRespectsTheLimit(t *testing.T) {
	draft := lesson.Draft{Title: "alpha beta gamma", Cue: "delta"}
	candidates := []lesson.Lesson{
		titled("a", "alpha beta", ""),
		titled("b", "alpha beta gamma", ""),
		titled("c", "beta gamma", ""),
	}

	assert.Equal(t, []string{"b"}, ids(Similar(draft, candidates, 1)))
}

func TestStopWordsDoNotMakeLessonsSimilar(t *testing.T) {
	draft := lesson.Draft{Title: "Run the tests before the release", Cue: "when it fails"}

	assert.Empty(t, Similar(draft, []lesson.Lesson{titled("a", "Run the linter before the commit", "when it fails")}, 5))
}
