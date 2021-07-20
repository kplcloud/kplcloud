/**
 * @Time : 2020/12/28 11:18 AM
 * @Author : solacowa@gmail.com
 * @File : redis_middleware
 * @Software: GoLand
 */

package redis

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis"
	redisclient "github.com/icowan/redis-client"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Middleware func(next redisclient.RedisClient) redisclient.RedisClient

type redisTraceServer struct {
	next   redisclient.RedisClient
	tracer opentracing.Tracer
}

func (s *redisTraceServer) Pipeline(ctx context.Context) redis.Pipeliner {
	panic("implement me")
}

func (s *redisTraceServer) Pipelined(ctx context.Context, fn func(pipeliner redis.Pipeliner) error) ([]redis.Cmder, error) {
	panic("implement me")
}

func (s *redisTraceServer) TxPipelined(ctx context.Context, fn func(pipeliner redis.Pipeliner) error) ([]redis.Cmder, error) {
	panic("implement me")
}

func (s *redisTraceServer) TxPipeline(ctx context.Context) redis.Pipeliner {
	panic("implement me")
}

func (s *redisTraceServer) Command(ctx context.Context) *redis.CommandsInfoCmd {
	panic("implement me")
}

func (s *redisTraceServer) ClientGetName(ctx context.Context) string {
	panic("implement me")
}

func (s *redisTraceServer) Echo(ctx context.Context, message interface{}) string {
	panic("implement me")
}

func (s *redisTraceServer) Quit(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) Unlink(ctx context.Context, keys ...string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) Dump(ctx context.Context, key string) string {
	panic("implement me")
}

func (s *redisTraceServer) Expire(ctx context.Context, key string, expiration time.Duration) bool {
	panic("implement me")
}

func (s *redisTraceServer) ExpireAt(ctx context.Context, key string, tm time.Time) bool {
	panic("implement me")
}

func (s *redisTraceServer) Migrate(ctx context.Context, host, port, key string, db int64, timeout time.Duration) error {
	panic("implement me")
}

func (s *redisTraceServer) Move(ctx context.Context, key string, db int64) bool {
	panic("implement me")
}

func (s *redisTraceServer) ObjectRefCount(ctx context.Context, key string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ObjectEncoding(ctx context.Context, key string) string {
	panic("implement me")
}

func (s *redisTraceServer) ObjectIdleTime(ctx context.Context, key string) time.Duration {
	panic("implement me")
}

func (s *redisTraceServer) Persist(ctx context.Context, key string) bool {
	panic("implement me")
}

func (s *redisTraceServer) PExpire(ctx context.Context, key string, expiration time.Duration) bool {
	panic("implement me")
}

func (s *redisTraceServer) PExpireAt(ctx context.Context, key string, tm time.Time) bool {
	panic("implement me")
}

func (s *redisTraceServer) PTTL(ctx context.Context, key string) time.Duration {
	panic("implement me")
}

func (s *redisTraceServer) RandomKey(ctx context.Context) string {
	panic("implement me")
}

func (s *redisTraceServer) Rename(ctx context.Context, key, newkey string) *redis.StatusCmd {
	panic("implement me")
}

func (s *redisTraceServer) RenameNX(ctx context.Context, key, newkey string) bool {
	panic("implement me")
}

func (s *redisTraceServer) Restore(ctx context.Context, key string, ttl time.Duration, value string) error {
	panic("implement me")
}

func (s *redisTraceServer) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) error {
	panic("implement me")
}

func (s *redisTraceServer) Sort(ctx context.Context, key string, sort *redis.Sort) []string {
	panic("implement me")
}

func (s *redisTraceServer) SortStore(ctx context.Context, key, store string, sort *redis.Sort) int64 {
	panic("implement me")
}

func (s *redisTraceServer) SortInterfaces(ctx context.Context, key string, sort *redis.Sort) []interface{} {
	panic("implement me")
}

func (s *redisTraceServer) Touch(ctx context.Context, keys ...string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) Type(ctx context.Context, key string) string {
	panic("implement me")
}

func (s *redisTraceServer) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64) {
	panic("implement me")
}

func (s *redisTraceServer) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64) {
	panic("implement me")
}

func (s *redisTraceServer) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64) {
	panic("implement me")
}

func (s *redisTraceServer) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64) {
	panic("implement me")
}

func (s *redisTraceServer) Append(ctx context.Context, key, value string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) int64 {
	panic("implement me")
}

func (s *redisTraceServer) BitOpAnd(ctx context.Context, destKey string, keys ...string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) BitOpOr(ctx context.Context, destKey string, keys ...string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) BitOpXor(ctx context.Context, destKey string, keys ...string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) BitOpNot(ctx context.Context, destKey string, key string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) BitPos(ctx context.Context, key string, bit int64, pos ...int64) int64 {
	panic("implement me")
}

func (s *redisTraceServer) Decr(ctx context.Context, key string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) DecrBy(ctx context.Context, key string, decrement int64) int64 {
	panic("implement me")
}

func (s *redisTraceServer) GetBit(ctx context.Context, key string, offset int64) int64 {
	panic("implement me")
}

func (s *redisTraceServer) GetRange(ctx context.Context, key string, start, end int64) string {
	panic("implement me")
}

func (s *redisTraceServer) GetSet(ctx context.Context, key string, value interface{}) string {
	panic("implement me")
}

func (s *redisTraceServer) IncrBy(ctx context.Context, key string, value int64) int64 {
	panic("implement me")
}

func (s *redisTraceServer) IncrByFloat(ctx context.Context, key string, value float64) float64 {
	panic("implement me")
}

func (s *redisTraceServer) MGet(ctx context.Context, keys ...string) []interface{} {
	panic("implement me")
}

func (s *redisTraceServer) MSet(ctx context.Context, pairs ...interface{}) error {
	panic("implement me")
}

func (s *redisTraceServer) MSetNX(ctx context.Context, pairs ...interface{}) (res bool) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MSetNX", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("pairs", pairs, "res", res)
		span.Finish()
	}()
	return s.next.MSetNX(ctx, pairs...)
}

func (s *redisTraceServer) SetBit(ctx context.Context, key string, offset int64, value int) int64 {
	panic("implement me")
}

func (s *redisTraceServer) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (res bool) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SetNX", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", key, "value", value, "expiration", expiration, "res", res)
		span.Finish()
	}()
	return s.next.SetNX(ctx, key, value, expiration)
}

func (s *redisTraceServer) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) (res bool) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SetXX", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", key, "value", value, "expiration", expiration, "res", res)
		span.Finish()
	}()
	return s.next.SetXX(ctx, key, value, expiration)
}

func (s *redisTraceServer) SetRange(ctx context.Context, key string, offset int64, value string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) StrLen(ctx context.Context, key string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) HExists(ctx context.Context, key, field string) bool {
	panic("implement me")
}

func (s *redisTraceServer) HIncrBy(ctx context.Context, key, field string, incr int64) int64 {
	panic("implement me")
}

func (s *redisTraceServer) HIncrByFloat(ctx context.Context, key, field string, incr float64) float64 {
	panic("implement me")
}

func (s *redisTraceServer) HKeys(ctx context.Context, key string) []string {
	panic("implement me")
}

func (s *redisTraceServer) HMGet(ctx context.Context, key string, fields ...string) []interface{} {
	panic("implement me")
}

func (s *redisTraceServer) HMSet(ctx context.Context, key string, fields map[string]interface{}) error {
	panic("implement me")
}

func (s *redisTraceServer) HSetNX(ctx context.Context, key, field string, value interface{}) bool {
	panic("implement me")
}

func (s *redisTraceServer) HVals(ctx context.Context, key string) []string {
	panic("implement me")
}

func (s *redisTraceServer) BLPop(ctx context.Context, timeout time.Duration, keys ...string) []string {
	panic("implement me")
}

func (s *redisTraceServer) BRPop(ctx context.Context, timeout time.Duration, keys ...string) []string {
	panic("implement me")
}

func (s *redisTraceServer) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) string {
	panic("implement me")
}

func (s *redisTraceServer) LIndex(ctx context.Context, key string, index int64) string {
	panic("implement me")
}

func (s *redisTraceServer) LInsert(ctx context.Context, key, op string, pivot, value interface{}) int64 {
	panic("implement me")
}

func (s *redisTraceServer) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) int64 {
	panic("implement me")
}

func (s *redisTraceServer) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) int64 {
	panic("implement me")
}

func (s *redisTraceServer) LPop(ctx context.Context, key string) string {
	panic("implement me")
}

func (s *redisTraceServer) LPushX(ctx context.Context, key string, value interface{}) int64 {
	panic("implement me")
}

func (s *redisTraceServer) LRange(ctx context.Context, key string, start, stop int64) []string {
	panic("implement me")
}

func (s *redisTraceServer) LRem(ctx context.Context, key string, count int64, value interface{}) int64 {
	panic("implement me")
}

func (s *redisTraceServer) LSet(ctx context.Context, key string, index int64, value interface{}) error {
	panic("implement me")
}

func (s *redisTraceServer) LTrim(ctx context.Context, key string, start, stop int64) error {
	panic("implement me")
}

func (s *redisTraceServer) RPopLPush(ctx context.Context, source, destination string) string {
	panic("implement me")
}

func (s *redisTraceServer) RPush(ctx context.Context, key string, values ...interface{}) int64 {
	panic("implement me")
}

func (s *redisTraceServer) RPushX(ctx context.Context, key string, value interface{}) int64 {
	panic("implement me")
}

func (s *redisTraceServer) SAdd(ctx context.Context, key string, members ...interface{}) int64 {
	panic("implement me")
}

func (s *redisTraceServer) SCard(ctx context.Context, key string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) SDiff(ctx context.Context, keys ...string) []string {
	panic("implement me")
}

func (s *redisTraceServer) SDiffStore(ctx context.Context, destination string, keys ...string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) SInter(ctx context.Context, keys ...string) []string {
	panic("implement me")
}

func (s *redisTraceServer) SInterStore(ctx context.Context, destination string, keys ...string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) SIsMember(ctx context.Context, key string, member interface{}) bool {
	panic("implement me")
}

func (s *redisTraceServer) SMembers(ctx context.Context, key string) []string {
	panic("implement me")
}

func (s *redisTraceServer) SMembersMap(ctx context.Context, key string) map[string]struct{} {
	panic("implement me")
}

func (s *redisTraceServer) SMove(ctx context.Context, source, destination string, member interface{}) bool {
	panic("implement me")
}

func (s *redisTraceServer) SPop(ctx context.Context, key string) string {
	panic("implement me")
}

func (s *redisTraceServer) SPopN(ctx context.Context, key string, count int64) []string {
	panic("implement me")
}

func (s *redisTraceServer) SRandMember(ctx context.Context, key string) string {
	panic("implement me")
}

func (s *redisTraceServer) SRandMemberN(ctx context.Context, key string, count int64) []string {
	panic("implement me")
}

func (s *redisTraceServer) SRem(ctx context.Context, key string, members ...interface{}) int64 {
	panic("implement me")
}

func (s *redisTraceServer) SUnion(ctx context.Context, keys ...string) []string {
	panic("implement me")
}

func (s *redisTraceServer) SUnionStore(ctx context.Context, destination string, keys ...string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) XAdd(ctx context.Context, a *redis.XAddArgs) string {
	panic("implement me")
}

func (s *redisTraceServer) XDel(ctx context.Context, stream string, ids ...string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) XLen(ctx context.Context, stream string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) XRange(ctx context.Context, stream, start, stop string) []redis.XMessage {
	panic("implement me")
}

func (s *redisTraceServer) XRangeN(ctx context.Context, stream, start, stop string, count int64) []redis.XMessage {
	panic("implement me")
}

func (s *redisTraceServer) XRevRange(ctx context.Context, stream string, start, stop string) []redis.XMessage {
	panic("implement me")
}

func (s *redisTraceServer) XRevRangeN(ctx context.Context, stream string, start, stop string, count int64) []redis.XMessage {
	panic("implement me")
}

func (s *redisTraceServer) XRead(ctx context.Context, a *redis.XReadArgs) []redis.XStream {
	panic("implement me")
}

func (s *redisTraceServer) XReadStreams(ctx context.Context, streams ...string) []redis.XStream {
	panic("implement me")
}

func (s *redisTraceServer) XGroupCreate(ctx context.Context, stream, group, start string) error {
	panic("implement me")
}

func (s *redisTraceServer) XGroupCreateMkStream(ctx context.Context, stream, group, start string) error {
	panic("implement me")
}

func (s *redisTraceServer) XGroupSetID(ctx context.Context, stream, group, start string) error {
	panic("implement me")
}

func (s *redisTraceServer) XGroupDestroy(ctx context.Context, stream, group string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) XGroupDelConsumer(ctx context.Context, stream, group, consumer string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) XReadGroup(ctx context.Context, a *redis.XReadGroupArgs) []redis.XStream {
	panic("implement me")
}

func (s *redisTraceServer) XAck(ctx context.Context, stream, group string, ids ...string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) XPending(stream, group string) *redis.XPending {
	panic("implement me")
}

func (s *redisTraceServer) XPendingExt(ctx context.Context, a *redis.XPendingExtArgs) []redis.XPendingExt {
	panic("implement me")
}

func (s *redisTraceServer) XClaim(ctx context.Context, a *redis.XClaimArgs) []redis.XMessage {
	panic("implement me")
}

func (s *redisTraceServer) XClaimJustID(ctx context.Context, a *redis.XClaimArgs) []string {
	panic("implement me")
}

func (s *redisTraceServer) XTrim(ctx context.Context, key string, maxLen int64) int64 {
	panic("implement me")
}

func (s *redisTraceServer) XTrimApprox(ctx context.Context, key string, maxLen int64) int64 {
	panic("implement me")
}

func (s *redisTraceServer) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) redis.ZWithKey {
	panic("implement me")
}

func (s *redisTraceServer) BZPopMin(timeout time.Duration, keys ...string) redis.ZWithKey {
	panic("implement me")
}

func (s *redisTraceServer) ZAddNX(ctx context.Context, key string, members ...redis.Z) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ZAddXX(ctx context.Context, key string, members ...redis.Z) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ZAddCh(ctx context.Context, key string, members ...redis.Z) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ZAddNXCh(ctx context.Context, key string, members ...redis.Z) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ZAddXXCh(ctx context.Context, key string, members ...redis.Z) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ZIncr(ctx context.Context, key string, member redis.Z) float64 {
	panic("implement me")
}

func (s *redisTraceServer) ZIncrNX(ctx context.Context, key string, member redis.Z) float64 {
	panic("implement me")
}

func (s *redisTraceServer) ZIncrXX(ctx context.Context, key string, member redis.Z) float64 {
	panic("implement me")
}

func (s *redisTraceServer) ZCount(ctx context.Context, key, min, max string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ZLexCount(ctx context.Context, key, min, max string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ZIncrBy(ctx context.Context, key string, increment float64, member string) float64 {
	panic("implement me")
}

func (s *redisTraceServer) ZInterStore(ctx context.Context, destination string, store redis.ZStore, keys ...string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ZPopMax(ctx context.Context, key string, count ...int64) []redis.Z {
	panic("implement me")
}

func (s *redisTraceServer) ZPopMin(ctx context.Context, key string, count ...int64) []redis.Z {
	panic("implement me")
}

func (s *redisTraceServer) ZRange(ctx context.Context, key string, start, stop int64) []string {
	panic("implement me")
}

func (s *redisTraceServer) ZRangeByScore(ctx context.Context, key string, opt redis.ZRangeBy) []string {
	panic("implement me")
}

func (s *redisTraceServer) ZRangeByLex(ctx context.Context, key string, opt redis.ZRangeBy) []string {
	panic("implement me")
}

func (s *redisTraceServer) ZRangeByScoreWithScores(ctx context.Context, key string, opt redis.ZRangeBy) []redis.Z {
	panic("implement me")
}

func (s *redisTraceServer) ZRank(ctx context.Context, key, member string) (int64, error) {
	panic("implement me")
}

func (s *redisTraceServer) ZRem(ctx context.Context, key string, members ...interface{}) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	panic("implement me")
}

func (s *redisTraceServer) ZRemRangeByScore(ctx context.Context, key, min, max string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ZRemRangeByLex(ctx context.Context, key, min, max string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ZRevRange(ctx context.Context, key string, start, stop int64) []string {
	panic("implement me")
}

func (s *redisTraceServer) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) []redis.Z {
	panic("implement me")
}

func (s *redisTraceServer) ZRevRangeByScore(ctx context.Context, key string, opt redis.ZRangeBy) []string {
	panic("implement me")
}

func (s *redisTraceServer) ZRevRangeByLex(ctx context.Context, key string, opt redis.ZRangeBy) []string {
	panic("implement me")
}

func (s *redisTraceServer) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt redis.ZRangeBy) []redis.Z {
	panic("implement me")
}

func (s *redisTraceServer) ZRevRank(ctx context.Context, key, member string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ZRevRank", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", key, "member", member)
		span.Finish()
	}()
	return s.next.ZRevRank(ctx, key, member)
}

func (s *redisTraceServer) ZScore(ctx context.Context, key, member string) float64 {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ZScore", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", key, "member", member)
		span.Finish()
	}()
	return s.next.ZScore(ctx, key, member)
}

func (s *redisTraceServer) ZUnionStore(ctx context.Context, dest string, store redis.ZStore, keys ...string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) PFAdd(ctx context.Context, key string, els ...interface{}) int64 {
	panic("implement me")
}

func (s *redisTraceServer) PFCount(ctx context.Context, keys ...string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) PFMerge(ctx context.Context, dest string, keys ...string) error {
	panic("implement me")
}

func (s *redisTraceServer) BgRewriteAOF(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) BgSave(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) ClientKill(ctx context.Context, ipPort string) error {
	panic("implement me")
}

func (s *redisTraceServer) ClientKillByFilter(ctx context.Context, keys ...string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ClientList(ctx context.Context) string {
	panic("implement me")
}

func (s *redisTraceServer) ClientPause(ctx context.Context, dur time.Duration) bool {
	panic("implement me")
}

func (s *redisTraceServer) ClientID(ctx context.Context) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ConfigGet(ctx context.Context, parameter string) []interface{} {
	panic("implement me")
}

func (s *redisTraceServer) ConfigResetStat(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) ConfigSet(ctx context.Context, parameter, value string) error {
	panic("implement me")
}

func (s *redisTraceServer) ConfigRewrite(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) DBSize(ctx context.Context) int64 {
	panic("implement me")
}

func (s *redisTraceServer) FlushAll(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) FlushAllAsync(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) FlushDB(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) FlushDBAsync(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) Info(ctx context.Context, section ...string) string {
	panic("implement me")
}

func (s *redisTraceServer) LastSave(ctx context.Context) int64 {
	panic("implement me")
}

func (s *redisTraceServer) Save(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) Shutdown(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) ShutdownSave(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) ShutdownNoSave(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) SlaveOf(ctx context.Context, host, port string) error {
	panic("implement me")
}

func (s *redisTraceServer) Time(ctx context.Context) time.Time {
	panic("implement me")
}

func (s *redisTraceServer) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	panic("implement me")
}

func (s *redisTraceServer) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	panic("implement me")
}

func (s *redisTraceServer) ScriptExists(ctx context.Context, hashes ...string) []bool {
	panic("implement me")
}

func (s *redisTraceServer) ScriptFlush(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) ScriptKill(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) ScriptLoad(ctx context.Context, script string) string {
	panic("implement me")
}

func (s *redisTraceServer) DebugObject(ctx context.Context, key string) string {
	panic("implement me")
}

func (s *redisTraceServer) PubSubChannels(ctx context.Context, pattern string) []string {
	panic("implement me")
}

func (s *redisTraceServer) PubSubNumSub(ctx context.Context, channels ...string) map[string]int64 {
	panic("implement me")
}

func (s *redisTraceServer) PubSubNumPat(ctx context.Context) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ClusterSlots(ctx context.Context) []redis.ClusterSlot {
	panic("implement me")
}

func (s *redisTraceServer) ClusterNodes(ctx context.Context) string {
	panic("implement me")
}

func (s *redisTraceServer) ClusterMeet(ctx context.Context, host, port string) error {
	panic("implement me")
}

func (s *redisTraceServer) ClusterForget(ctx context.Context, nodeID string) error {
	panic("implement me")
}

func (s *redisTraceServer) ClusterReplicate(ctx context.Context, nodeID string) error {
	panic("implement me")
}

func (s *redisTraceServer) ClusterResetSoft(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) ClusterResetHard(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) ClusterInfo(ctx context.Context) string {
	panic("implement me")
}

func (s *redisTraceServer) ClusterKeySlot(ctx context.Context, key string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ClusterGetKeysInSlot(ctx context.Context, slot int, count int) []string {
	panic("implement me")
}

func (s *redisTraceServer) ClusterCountFailureReports(ctx context.Context, nodeID string) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ClusterCountKeysInSlot(ctx context.Context, slot int) int64 {
	panic("implement me")
}

func (s *redisTraceServer) ClusterDelSlots(ctx context.Context, slots ...int) error {
	panic("implement me")
}

func (s *redisTraceServer) ClusterDelSlotsRange(ctx context.Context, min, max int) error {
	panic("implement me")
}

func (s *redisTraceServer) ClusterSaveConfig(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) ClusterSlaves(ctx context.Context, nodeID string) []string {
	panic("implement me")
}

func (s *redisTraceServer) ClusterFailover(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) ClusterAddSlots(ctx context.Context, slots ...int) error {
	panic("implement me")
}

func (s *redisTraceServer) ClusterAddSlotsRange(ctx context.Context, min, max int) error {
	panic("implement me")
}

func (s *redisTraceServer) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) int64 {
	panic("implement me")
}

func (s *redisTraceServer) GeoPos(ctx context.Context, key string, members ...string) []*redis.GeoPos {
	panic("implement me")
}

func (s *redisTraceServer) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) []redis.GeoLocation {
	panic("implement me")
}

func (s *redisTraceServer) GeoRadiusRO(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) []redis.GeoLocation {
	panic("implement me")
}

func (s *redisTraceServer) GeoRadiusByMember(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) []redis.GeoLocation {
	panic("implement me")
}

func (s *redisTraceServer) GeoRadiusByMemberRO(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) []redis.GeoLocation {
	panic("implement me")
}

func (s *redisTraceServer) GeoDist(ctx context.Context, key string, member1, member2, unit string) float64 {
	panic("implement me")
}

func (s *redisTraceServer) GeoHash(ctx context.Context, key string, members ...string) []string {
	panic("implement me")
}

func (s *redisTraceServer) ReadOnly(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) ReadWrite(ctx context.Context) error {
	panic("implement me")
}

func (s *redisTraceServer) MemoryUsage(ctx context.Context, key string, samples ...int) int64 {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MemoryUsage", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", key, "samples", samples)
		span.Finish()
	}()
	return s.next.MemoryUsage(ctx, key, samples...)
}

func (s *redisTraceServer) Close(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Close", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.Finish()
	}()
	return s.next.Close(ctx)
}

func (s *redisTraceServer) Del(ctx context.Context, k string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Del", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", k, "error", err)
		span.Finish()
	}()
	return s.next.Del(ctx, k)
}

func (s *redisTraceServer) Exists(ctx context.Context, keys ...string) int64 {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Exists", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("keys", strings.Join(keys, ","))
		span.Finish()
	}()
	return s.next.Exists(ctx, keys...)
}

func (s *redisTraceServer) HDel(ctx context.Context, k string, field string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "HDel", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", k, "field", field, "error", err)
		span.Finish()
	}()
	return s.next.HDel(ctx, k, field)
}

func (s *redisTraceServer) HDelAll(ctx context.Context, k string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "HDelAll", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", k, "error", err)
		span.Finish()
	}()
	return s.next.HDelAll(ctx, k)
}

func (s *redisTraceServer) HGet(ctx context.Context, k string, field string) (res string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "HGet", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", k, "field", field, "error", err)
		span.Finish()
	}()
	return s.next.HGet(ctx, k, field)
}

func (s *redisTraceServer) HGetAll(ctx context.Context, k string) (res map[string]string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "HGetAll", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", k, "error", err)
		span.Finish()
	}()
	return s.next.HGetAll(ctx, k)
}

func (s *redisTraceServer) HLen(ctx context.Context, k string) (res int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "HLen", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", k, "res", res, "error", err)
		span.Finish()
	}()
	return s.next.HLen(ctx, k)
}

func (s *redisTraceServer) HSet(ctx context.Context, k string, field string, v interface{}) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "HSet", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", k, "field", field, "v", field, "error", err)
		span.Finish()
	}()
	return s.next.HSet(ctx, k, field, v)
}

func (s *redisTraceServer) Incr(ctx context.Context, key string, exp time.Duration) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Incr", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", key, "key", exp.String())
		span.Finish()
	}()
	return s.next.Incr(ctx, key, exp)
}

func (s *redisTraceServer) Keys(ctx context.Context, pattern string) (res []string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Keys", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("pattern", pattern, "resLen", len(res))
		span.Finish()
	}()
	return s.next.Keys(ctx, pattern)
}

func (s *redisTraceServer) LLen(ctx context.Context, key string) int64 {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "LLen", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", key)
		span.Finish()
	}()
	return s.next.LLen(ctx, key)
}

func (s *redisTraceServer) LPush(ctx context.Context, key string, val interface{}) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "LPush", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", key, "val", val)
		span.Finish()
	}()
	return s.next.LPush(ctx, key, val)
}

func (s *redisTraceServer) Ping(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Ping", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		//span.LogKV("ping", "ping")
		span.Finish()
	}()
	return s.next.Ping(ctx)
}

func (s *redisTraceServer) Publish(ctx context.Context, channel string, message interface{}) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Publish", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("channel", channel)
		span.Finish()
	}()
	return s.next.Publish(ctx, channel, message)
}

func (s *redisTraceServer) RPop(ctx context.Context, key string) (res string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "RPop", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", key)
		span.Finish()
	}()
	return s.next.RPop(ctx, key)
}

func (s *redisTraceServer) Set(ctx context.Context, k string, v interface{}, expir ...time.Duration) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Set", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", k, "expir", expir)
		span.Finish()
	}()
	return s.next.Set(ctx, k, v, expir...)
}

func (s *redisTraceServer) SetPrefix(ctx context.Context, prefix string) redisclient.RedisClient {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Subscribe", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("prefix", prefix)
		span.Finish()
	}()
	return s.next.SetPrefix(ctx, prefix)
}

func (s *redisTraceServer) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Subscribe", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", strings.Join(channels, ","))
		span.Finish()
	}()
	return s.next.Subscribe(ctx, channels...)
}

func (s *redisTraceServer) TTL(ctx context.Context, key string) (t time.Duration) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "TypeOf", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", key, "t", t.String())
		span.Finish()
	}()
	return s.next.TTL(ctx, key)
}

func (s *redisTraceServer) TypeOf(ctx context.Context, key string) (res string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "TypeOf", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", key)
		span.Finish()
	}()
	return s.next.TypeOf(ctx, key)
}

func (s *redisTraceServer) ZAdd(ctx context.Context, k string, score float64, member interface{}) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ZAdd", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", k, "score", score, "member", member)
		span.Finish()
	}()
	return s.next.ZAdd(ctx, k, score, member)
}

func (s *redisTraceServer) ZCard(ctx context.Context, k string) (res int64, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ZCard", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", k)
		span.Finish()
	}()
	return s.next.ZCard(ctx, k)
}

func (s *redisTraceServer) ZRangeWithScores(ctx context.Context, k string, start, stop int64) (res []redis.Z, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ZRangeWithScores", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", k, "start", start, "stop", stop)
		span.Finish()
	}()
	return s.next.ZRangeWithScores(ctx, k, start, stop)
}

func (s *redisTraceServer) Get(ctx context.Context, k string) (v string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Get", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "Redis",
	})
	defer func() {
		span.LogKV("key", k, "v", v, "error", err)
		span.Finish()
	}()
	return s.next.Get(ctx, k)
}

func NewRedisMiddleware(svc redisclient.RedisClient, tracer opentracing.Tracer) Middleware {
	return func(next redisclient.RedisClient) redisclient.RedisClient {
		return &redisTraceServer{
			next:   svc,
			tracer: tracer,
		}
	}
}
