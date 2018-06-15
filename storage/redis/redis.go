package redis

import (
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type Redis struct {
	cli    *redis.Client
	prefix string

	mut   sync.RWMutex
	index map[string]int32
	lits  map[string][]string
}

func NewRedis(cli *redis.Client) *Redis {
	return &Redis{
		cli:    cli,
		prefix: "pool",
		index:  map[string]int32{},
		lits:   map[string][]string{},
	}
}

func (r *Redis) IsExist(prefix, raw string) (bool, error) {
	k := strings.Join([]string{r.prefix, prefix, "update"}, ":")
	return r.cli.HExists(k, raw).Result()
}

func (r *Redis) update(prefix, raw string) error {
	k := strings.Join([]string{r.prefix, prefix, "update"}, ":")
	data := time.Now().Format(time.RFC3339)
	_, err := r.cli.HSet(k, raw, data).Result()
	return err
}

func (r *Redis) IncUsed(prefix, raw string) error {
	k := strings.Join([]string{r.prefix, prefix, "used"}, ":")
	_, err := r.cli.HIncrBy(k, raw, 1).Result()
	if err != nil {
		return err
	}
	return r.update(prefix, raw)
}

func (r *Redis) IncCheck(prefix, raw string) error {
	k := strings.Join([]string{r.prefix, prefix, "check"}, ":")
	_, err := r.cli.HIncrBy(k, raw, 1).Result()
	if err != nil {
		return err
	}
	data := time.Now().Format(time.RFC3339)
	k = strings.Join([]string{r.prefix, prefix, "ready"}, ":")
	_, err = r.cli.HSet(k, raw, data).Result()
	return r.update(prefix, raw)
}

func (r *Redis) IncFailure(prefix, raw string) error {
	k := strings.Join([]string{r.prefix, prefix, "failure"}, ":")
	_, err := r.cli.HIncrBy(k, raw, 1).Result()
	if err != nil {
		return err
	}
	k = strings.Join([]string{r.prefix, prefix, "ready"}, ":")
	r.cli.HDel(k, raw)
	return r.update(prefix, raw)
}

func (r *Redis) GetRandom(prefix string) (string, error) {
	r.mut.Lock()
	defer r.mut.Unlock()
	tmpi := r.index[prefix]
	tmpl := r.lits[prefix]
	if index := tmpi; int(index) != len(tmpl) {
		r.index[prefix] = tmpi + 1
		return tmpl[index], nil
	}
	k := strings.Join([]string{r.prefix, prefix, "ready"}, ":")
	m, err := r.cli.HKeys(k).Result()
	if err != nil {
		return "", err
	}

	r.lits[prefix] = m
	r.index[prefix] = 0
	if len(m) == 0 {
		return "", nil
	}
	r.index[prefix] = 1
	return m[0], nil
}

func (r *Redis) GetAll(prefix string) ([]string, error) {
	r.mut.Lock()
	defer r.mut.Unlock()

	k := strings.Join([]string{r.prefix, prefix, "ready"}, ":")
	m, err := r.cli.HKeys(k).Result()
	if err != nil {
		return nil, err
	}

	return m, nil
}
