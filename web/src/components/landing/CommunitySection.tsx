import { motion } from 'framer-motion'
import AnimatedSection from './AnimatedSection'

interface CardProps {
  quote: string
  authorName: string
  handle: string
  avatarUrl: string
  tweetUrl: string
  delay: number
}

function TestimonialCard({ quote, authorName, delay }: CardProps) {
  return (
    <motion.div
      className="p-6 rounded-xl"
      style={{
        background: 'var(--brand-dark-gray)',
        border: '1px solid rgba(240, 185, 11, 0.1)',
      }}
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      transition={{ delay }}
      whileHover={{ scale: 1.05 }}
    >
      <p className="text-lg mb-4" style={{ color: 'var(--brand-light-gray)' }}>
        "{quote}"
      </p>
      <div className="flex items-center gap-2">
        <div
          className="w-8 h-8 rounded-full"
          style={{ background: 'var(--binance-yellow)' }}
        />
        <span
          className="text-sm font-semibold"
          style={{ color: 'var(--text-secondary)' }}
        >
          {authorName}
        </span>
      </div>
    </motion.div>
  )
}

export default function CommunitySection() {
  const staggerContainer = {
    animate: { transition: { staggerChildren: 0.1 } },
  }

  // 推特内容整合（保持原三列布局，超出自动换行）
  const items: CardProps[] = [
    {
      quote:
        ' AI 量化交易系统  DeepSeek、Qwen 等大语言模型，打造的通用架构 AI 交易操作系统，完成了从决策、到交易、再到复盘的闭环。',
      authorName: '乔布斯',
      handle: '@乔布斯',
      avatarUrl:
        '#',
      tweetUrl:
        '#',
      delay: 0,
    },
    {
      quote:
        '跑了一晚上 AI 自动交易，太有意思了，就看 AI 在那一会开空一会开多，一顿操作，虽然看不懂为什么，但是一晚上帮我赚了 6% 收益',
      authorName: '哈瓦斯',
      handle: '@哈瓦斯',
      avatarUrl:
        '#',
      tweetUrl: '#',
      delay: 0.1,
    },
    {
      quote:
        '使用 的 AI 交易系统已经有一段时间了，效果真的很惊人。AI 能够根据市场变化自动调整策略，帮助我在波动的市场中保持盈利。强烈推荐给所有对加密交易感兴趣的人！',
      authorName: 'SU7',
      handle: '@SU7',
      avatarUrl:
        '#',
      tweetUrl: '#',
      delay: 0.15,
    },
  ]

  return (
    <AnimatedSection>
      <div className="max-w-7xl mx-auto">
        <motion.div
          className="grid md:grid-cols-3 gap-6"
          variants={staggerContainer}
          initial="initial"
          whileInView="animate"
          viewport={{ once: true }}
        >
          {items.map((item, idx) => (
            <TestimonialCard key={idx} {...item} />
          ))}
        </motion.div>
      </div>
    </AnimatedSection>
  )
}
